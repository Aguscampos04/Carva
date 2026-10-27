package integration_test

import (
	"context"
	"encoding/json"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"net/http"
	"strconv"
	"testing"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	backlogdomain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestSprintBacklogAssignmentPreservesPlanningSnapshot(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(context.Background(), "Carva")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := appbacklog.NewService(backlogRepository)
	provisionalPoints := 5
	item, err := backlogService.Create(context.Background(), projectValue.ID, backlogdomain.Input{EstimatedHours: hours.Ptr(800),
		Title: "Plan inicial", Description: "Descripción original", Priority: backlogdomain.PriorityHigh,
		ProvisionalStoryPoints: &provisionalPoints, AcceptanceCriteria: []string{"Criterio original"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), backlogRepository)
	sprintValue, err := sprintService.Create(context.Background(), projectValue.ID, "Entregar el MVP")
	if err != nil {
		t.Fatal(err)
	}
	handler := httpapi.NewHandler(projectService, backlogService, sprintService)
	basePath := "/api/projects/" + projectValue.ID + "/sprints/" + strconv.FormatInt(sprintValue.ID, 10) + "/backlog-items"

	assigned := performJSONRequest(handler, http.MethodPost, basePath, map[string]any{"backlog_item_id": item.ID})
	if assigned.Code != http.StatusCreated {
		t.Fatalf("assign status = %d, body = %s; want %d", assigned.Code, assigned.Body, http.StatusCreated)
	}

	updated := performJSONRequest(handler, http.MethodPut, "/api/projects/"+projectValue.ID+"/backlog/"+item.ID, map[string]any{
		"title": "Título actualizado", "description": "Descripción actualizada", "priority": "Baja",
		"provisional_story_points": 8, "estimated_hours": 8, "acceptance_criteria": []string{"Criterio actualizado"},
	})
	if updated.Code != http.StatusOK {
		t.Fatalf("update backlog status = %d, body = %s; want %d", updated.Code, updated.Body, http.StatusOK)
	}

	planned := performJSONRequest(handler, http.MethodGet, basePath, nil)
	if planned.Code != http.StatusOK {
		t.Fatalf("list planned stories status = %d, body = %s; want %d", planned.Code, planned.Body, http.StatusOK)
	}
	var snapshots []struct {
		BacklogItemID string   `json:"backlog_item_id"`
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Priority      string   `json:"priority"`
		Status        string   `json:"status"`
		Provisional   *int     `json:"provisional_story_points"`
		StoryPoints   *int     `json:"story_points"`
		Criteria      []string `json:"acceptance_criteria"`
	}
	if err := json.NewDecoder(planned.Body).Decode(&snapshots); err != nil {
		t.Fatalf("decode planned stories: %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("planned snapshots = %d, want 1", len(snapshots))
	}
	snapshot := snapshots[0]
	if snapshot.BacklogItemID != item.ID || snapshot.Title != "Plan inicial" || snapshot.Description != "Descripción original" ||
		snapshot.Priority != "Alta" || snapshot.Status != string(backlogdomain.Pending) || snapshot.Provisional == nil || *snapshot.Provisional != provisionalPoints ||
		snapshot.StoryPoints != nil || len(snapshot.Criteria) != 1 || snapshot.Criteria[0] != "Criterio original" {
		t.Fatalf("planned snapshot = %+v; want original Backlog values", snapshot)
	}
}

func TestSprintBacklogAssignmentRejectsForeignAndDuplicateStories(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	firstProject, err := projectService.Create(context.Background(), "First")
	if err != nil {
		t.Fatal(err)
	}
	secondProject, err := projectService.Create(context.Background(), "Second")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := appbacklog.NewService(backlogRepository)
	item, err := backlogService.Create(context.Background(), secondProject.ID, backlogdomain.Input{EstimatedHours: hours.Ptr(800),
		Title: "Foreign story", Description: "Description", Priority: backlogdomain.PriorityMedium,
		AcceptanceCriteria: []string{"Criterion"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), backlogRepository)
	sprintValue, err := sprintService.Create(context.Background(), firstProject.ID, "Goal")
	if err != nil {
		t.Fatal(err)
	}
	handler := httpapi.NewHandler(projectService, backlogService, sprintService)
	path := "/api/projects/" + firstProject.ID + "/sprints/" + strconv.FormatInt(sprintValue.ID, 10) + "/backlog-items"
	body := map[string]any{"backlog_item_id": item.ID}
	foreign := performJSONRequest(handler, http.MethodPost, path, body)
	if foreign.Code != http.StatusNotFound {
		t.Fatalf("foreign item assignment status = %d, want %d", foreign.Code, http.StatusNotFound)
	}

	localItem, err := backlogService.Create(context.Background(), firstProject.ID, backlogdomain.Input{EstimatedHours: hours.Ptr(800),
		Title: "Local story", Description: "Description", Priority: backlogdomain.PriorityMedium,
		AcceptanceCriteria: []string{"Criterion"},
	})
	if err != nil {
		t.Fatal(err)
	}
	body["backlog_item_id"] = localItem.ID
	first := performJSONRequest(handler, http.MethodPost, path, body)
	if first.Code != http.StatusCreated {
		t.Fatalf("first assignment status = %d, body = %s", first.Code, first.Body)
	}
	duplicate := performJSONRequest(handler, http.MethodPost, path, body)
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate assignment status = %d, want %d", duplicate.Code, http.StatusConflict)
	}
	planned := performJSONRequest(handler, http.MethodGet, path, nil)
	var snapshots []json.RawMessage
	if err := json.NewDecoder(planned.Body).Decode(&snapshots); err != nil {
		t.Fatal(err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("planned snapshots = %d, want exactly one", len(snapshots))
	}
}

func TestCompleteStoryRequiresVerificationAndUpdatesSprintAndBacklog(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(context.Background(), "Carva")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := appbacklog.NewService(backlogRepository)
	item, err := backlogService.Create(context.Background(), projectValue.ID, backlogdomain.Input{EstimatedHours: hours.Ptr(800),
		Title: "Historia", Description: "Descripción", Priority: backlogdomain.PriorityHigh,
		AcceptanceCriteria: []string{"Criterio verificado"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), backlogRepository)
	sprintValue, err := sprintService.Create(context.Background(), projectValue.ID, "Objetivo")
	if err != nil {
		t.Fatal(err)
	}
	handler := httpapi.NewHandler(projectService, backlogService, sprintService)
	basePath := "/api/projects/" + projectValue.ID + "/sprints/" + strconv.FormatInt(sprintValue.ID, 10) + "/backlog-items"
	if got := performJSONRequest(handler, http.MethodPost, basePath, map[string]any{"backlog_item_id": item.ID}); got.Code != http.StatusCreated {
		t.Fatalf("assign story status = %d, body = %s", got.Code, got.Body)
	}
	completionPath := basePath + "/" + item.ID + "/completion"

	missing := performJSONRequest(handler, http.MethodPost, completionPath, map[string]any{})
	if missing.Code != http.StatusBadRequest {
		t.Fatalf("missing verification status = %d, want %d", missing.Code, http.StatusBadRequest)
	}
	unknown := performJSONRequest(handler, http.MethodPost, completionPath, map[string]any{"acceptance_criteria_verified": true, "verified": true})
	if unknown.Code != http.StatusBadRequest {
		t.Fatalf("unknown completion field status = %d, want %d", unknown.Code, http.StatusBadRequest)
	}
	unverified := performJSONRequest(handler, http.MethodPost, completionPath, map[string]any{"acceptance_criteria_verified": false})
	if unverified.Code != http.StatusConflict {
		t.Fatalf("unverified completion status = %d, want %d", unverified.Code, http.StatusConflict)
	}
	backlogBefore := performJSONRequest(handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/backlog", nil)
	var backlogBeforeItems []backlogdomain.Item
	if err := json.NewDecoder(backlogBefore.Body).Decode(&backlogBeforeItems); err != nil {
		t.Fatal(err)
	}
	if backlogBeforeItems[0].Status != backlogdomain.Pending {
		t.Fatalf("Backlog status after rejected completion = %q, want %q", backlogBeforeItems[0].Status, backlogdomain.Pending)
	}

	completed := performJSONRequest(handler, http.MethodPost, completionPath, map[string]any{"acceptance_criteria_verified": true})
	if completed.Code != http.StatusOK {
		t.Fatalf("complete story status = %d, body = %s", completed.Code, completed.Body)
	}
	retried := performJSONRequest(handler, http.MethodPost, completionPath, map[string]any{"acceptance_criteria_verified": true})
	if retried.Code != http.StatusOK {
		t.Fatalf("retry completed story status = %d, body = %s", retried.Code, retried.Body)
	}
	var sprintStory struct {
		Status        string `json:"status"`
		PlannedStatus string `json:"planned_status"`
		Verified      bool   `json:"acceptance_criteria_verified"`
		Title         string `json:"title"`
	}
	if err := json.NewDecoder(completed.Body).Decode(&sprintStory); err != nil {
		t.Fatal(err)
	}
	if sprintStory.Status != string(backlogdomain.Completed) || sprintStory.PlannedStatus != string(backlogdomain.Pending) || !sprintStory.Verified || sprintStory.Title != "Historia" {
		t.Fatalf("completion response = %+v, want current completed status and original plan", sprintStory)
	}

	planned := performJSONRequest(handler, http.MethodGet, basePath, nil)
	var plannedItems []struct {
		Status        string `json:"status"`
		PlannedStatus string `json:"planned_status"`
		Verified      bool   `json:"acceptance_criteria_verified"`
		Title         string `json:"title"`
	}
	if err := json.NewDecoder(planned.Body).Decode(&plannedItems); err != nil {
		t.Fatal(err)
	}
	if len(plannedItems) != 1 || plannedItems[0].Status != string(backlogdomain.Completed) || plannedItems[0].PlannedStatus != string(backlogdomain.Pending) || !plannedItems[0].Verified || plannedItems[0].Title != "Historia" {
		t.Fatalf("Sprint planning query = %+v, want completed execution status with original plan", plannedItems)
	}
	backlogAfter := performJSONRequest(handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/backlog", nil)
	var backlogItems []backlogdomain.Item
	if err := json.NewDecoder(backlogAfter.Body).Decode(&backlogItems); err != nil {
		t.Fatal(err)
	}
	if backlogItems[0].Status != backlogdomain.Completed {
		t.Fatalf("Backlog status = %q, want %q", backlogItems[0].Status, backlogdomain.Completed)
	}
}

func TestCompleteStoryRequiresSprintAssignmentAndRollsBackOnFailure(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(context.Background(), "Carva")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := appbacklog.NewService(backlogRepository)
	item, err := backlogService.Create(context.Background(), projectValue.ID, backlogdomain.Input{EstimatedHours: hours.Ptr(800),
		Title: "Historia", Description: "Descripción", Priority: backlogdomain.PriorityHigh,
		AcceptanceCriteria: []string{"Criterio"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), backlogRepository)
	sprintValue, err := sprintService.Create(context.Background(), projectValue.ID, "Objetivo")
	if err != nil {
		t.Fatal(err)
	}
	handler := httpapi.NewHandler(projectService, backlogService, sprintService)
	basePath := "/api/projects/" + projectValue.ID + "/sprints/" + strconv.FormatInt(sprintValue.ID, 10) + "/backlog-items/" + item.ID + "/completion"

	unassigned := performJSONRequest(handler, http.MethodPost, basePath, map[string]any{"acceptance_criteria_verified": true})
	if unassigned.Code != http.StatusNotFound {
		t.Fatalf("unassigned story completion status = %d, want %d", unassigned.Code, http.StatusNotFound)
	}
	backlogItems, err := backlogRepository.ListByProject(context.Background(), projectValue.ID)
	if err != nil || backlogItems[0].Status != backlogdomain.Pending {
		t.Fatalf("Backlog after unassigned completion = %+v, error = %v", backlogItems, err)
	}

	assignmentPath := "/api/projects/" + projectValue.ID + "/sprints/" + strconv.FormatInt(sprintValue.ID, 10) + "/backlog-items"
	if assigned := performJSONRequest(handler, http.MethodPost, assignmentPath, map[string]any{"backlog_item_id": item.ID}); assigned.Code != http.StatusCreated {
		t.Fatalf("assign story status = %d, body = %s", assigned.Code, assigned.Body)
	}
	if _, err := db.Exec(`CREATE TRIGGER reject_sprint_completion BEFORE UPDATE OF completion_verified ON sprint_backlog_items BEGIN SELECT RAISE(ABORT, 'forced failure'); END`); err != nil {
		t.Fatal(err)
	}
	failed := performJSONRequest(handler, http.MethodPost, basePath, map[string]any{"acceptance_criteria_verified": true})
	if failed.Code != http.StatusInternalServerError {
		t.Fatalf("forced persistence failure status = %d, body = %s", failed.Code, failed.Body)
	}
	backlogItems, err = backlogRepository.ListByProject(context.Background(), projectValue.ID)
	if err != nil || backlogItems[0].Status != backlogdomain.Pending {
		t.Fatalf("Backlog after transaction rollback = %+v, error = %v", backlogItems, err)
	}
}
