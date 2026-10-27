package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	appmetrics "github.com/Aguscampos04/Carva/internal/application/metrics"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	domain "github.com/Aguscampos04/Carva/internal/domain/metrics"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestUS20StoryPointMetricsUseSprintSnapshots(t *testing.T) {
	ctx := context.Background()
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "Métricas")
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), nil)
	createSprint := func(goal string) int64 {
		t.Helper()
		value, err := sprintService.Create(ctx, projectValue.ID, goal)
		if err != nil {
			t.Fatal(err)
		}
		return value.ID
	}
	first := createSprint("Primero")
	second := createSprint("Segundo")
	third := createSprint("Abierto")

	insertSnapshot := func(sprintID int64, itemID string, storyPoints any, completed bool) {
		t.Helper()
		_, err := db.ExecContext(ctx, `INSERT INTO sprint_backlog_items (
			sprint_id, backlog_item_id, title, description, priority, status,
			provisional_story_points, story_points, acceptance_criteria, completion_verified
		) VALUES (?, ?, ?, ?, ?, ?, NULL, ?, '[]', ?)`,
			sprintID, itemID, itemID, "Description", "Media", "Pendiente", storyPoints, completed)
		if err != nil {
			t.Fatal(err)
		}
	}
	insertSnapshot(first, "a", 3, true)
	insertSnapshot(first, "b", 5, false)
	insertSnapshot(second, "b", 5, true)
	insertSnapshot(second, "c", 2, false)
	insertSnapshot(third, "d", 8, false)
	if _, err := db.ExecContext(ctx, `UPDATE sprints SET closed = 1 WHERE id IN (?, ?)`, first, second); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO backlog_items (id, project_id, title, description, priority, status, story_points, acceptance_criteria)
		VALUES ('b', ?, 'Current title', 'Description', 'Media', 'Pendiente', 99, '[]')`, projectValue.ID); err != nil {
		t.Fatal(err)
	}

	handler := httpapi.WithMetrics(
		httpapi.NewHandler(projectService, nil, sprintService),
		appmetrics.NewService(sqlite.NewMetricsRepository(db)),
	)
	response := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/metrics/story-points", "")
	if response.Code != http.StatusOK {
		t.Fatalf("project metrics status = %d, body = %s", response.Code, response.Body)
	}
	var projectResult struct {
		Metrics domain.StoryPointMetrics `json:"metrics"`
		Sprints []struct {
			ID      int64                    `json:"id"`
			Closed  bool                     `json:"closed"`
			Metrics domain.StoryPointMetrics `json:"metrics"`
		} `json:"sprints"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &projectResult); err != nil {
		t.Fatal(err)
	}
	if len(projectResult.Sprints) != 3 || *projectResult.Metrics.PlannedStoryPoints != 23 || *projectResult.Metrics.CompletedStoryPoints != 8 || *projectResult.Metrics.TeamVelocity != 4 || *projectResult.Metrics.CompletionPercentage != 40 {
		t.Fatalf("project metrics = %+v", projectResult)
	}
	if projectResult.Sprints[0].ID != first || !projectResult.Sprints[0].Closed || *projectResult.Sprints[0].Metrics.PlannedStoryPoints != 8 || *projectResult.Sprints[0].Metrics.CompletedStoryPoints != 3 || *projectResult.Sprints[0].Metrics.TeamVelocity != 3 || *projectResult.Sprints[0].Metrics.CompletionPercentage != 50 {
		t.Fatalf("first Sprint metrics = %+v", projectResult.Sprints[0])
	}
	if projectResult.Sprints[2].Metrics.TeamVelocity != nil {
		t.Fatalf("open Sprint velocity = %v, want unavailable", projectResult.Sprints[2].Metrics.TeamVelocity)
	}

	sprintResponse := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/sprints/"+itoa(first)+"/metrics/story-points", "")
	if sprintResponse.Code != http.StatusOK {
		t.Fatalf("Sprint metrics status = %d, body = %s", sprintResponse.Code, sprintResponse.Body)
	}
	for _, target := range []string{
		"/api/projects/missing/metrics/story-points",
		"/api/projects/" + projectValue.ID + "/sprints/999/metrics/story-points",
	} {
		if got := performJSON(t, handler, http.MethodGet, target, "").Code; got != http.StatusNotFound {
			t.Fatalf("GET %s status = %d, want 404", target, got)
		}
	}
	if got := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/sprints/nope/metrics/story-points", "").Code; got != http.StatusBadRequest {
		t.Fatalf("invalid Sprint ID status = %d, want 400", got)
	}
}

func TestUS20MissingPointsAreNotReportedAsZero(t *testing.T) {
	ctx := context.Background()
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "Sin puntos")
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), nil)
	value, err := sprintService.Create(ctx, projectValue.ID, "Incomplete estimation")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO sprint_backlog_items (
		sprint_id, backlog_item_id, title, description, priority, status, acceptance_criteria, completion_verified
	) VALUES (?, 'missing', 'Missing points', 'Description', 'Media', 'Pendiente', '[]', 0)`, value.ID); err != nil {
		t.Fatal(err)
	}
	handler := httpapi.WithMetrics(httpapi.NewHandler(projectService, nil, sprintService), appmetrics.NewService(sqlite.NewMetricsRepository(db)))
	response := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/sprints/"+itoa(value.ID)+"/metrics/story-points", "")
	var result domain.StoryPointMetrics
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &result) != nil {
		t.Fatalf("metrics = %d %s", response.Code, response.Body)
	}
	if result.PlannedStoryPoints != nil || result.CompletedStoryPoints == nil || *result.CompletedStoryPoints != 0 || result.TeamVelocity != nil || result.CompletionPercentage == nil || *result.CompletionPercentage != 0 {
		t.Fatalf("metrics with missing points = %+v", result)
	}
}
