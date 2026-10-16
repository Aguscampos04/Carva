package integration_test

import (
	"context"
	"encoding/json"
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
	item, err := backlogService.Create(context.Background(), projectValue.ID, backlogdomain.Input{
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
		"provisional_story_points": 8, "acceptance_criteria": []string{"Criterio actualizado"},
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
	item, err := backlogService.Create(context.Background(), secondProject.ID, backlogdomain.Input{
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

	localItem, err := backlogService.Create(context.Background(), firstProject.ID, backlogdomain.Input{
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
