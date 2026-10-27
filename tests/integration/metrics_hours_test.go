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

func TestUS21HoursMetricsUseSprintAssociationAndDistinctProjectStories(t *testing.T) {
	ctx := context.Background()
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "Métricas de horas")
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
	first, second := createSprint("Primero"), createSprint("Segundo")
	if _, err := db.ExecContext(ctx, `INSERT INTO members(id,project_id,name) VALUES('member',?,'Ana')`, projectValue.ID); err != nil {
		t.Fatal(err)
	}
	insertStory := func(id string, estimate any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, `INSERT INTO backlog_items (
			id,project_id,title,description,priority,status,acceptance_criteria,estimated_hours_hundredths
		) VALUES (?, ?, ?, 'Descripción', 'Media', 'Pendiente', '[]', ?)`, id, projectValue.ID, id, estimate); err != nil {
			t.Fatal(err)
		}
	}
	insertStory("a", 400)
	insertStory("b", 300)
	insertSnapshot := func(sprintID int64, itemID string, estimate int64) {
		t.Helper()
		_, err := db.ExecContext(ctx, `INSERT INTO sprint_backlog_items (
			sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria,estimated_hours_hundredths
		) VALUES (?, ?, ?, 'Descripción', 'Media', 'Pendiente', '[]', ?)`, sprintID, itemID, itemID, estimate)
		if err != nil {
			t.Fatal(err)
		}
	}
	insertSnapshot(first, "a", 400)
	insertSnapshot(first, "b", 300)
	insertSnapshot(second, "a", 400)
	insertEffort := func(id string, hours int64, sprintID any, deleted any) {
		t.Helper()
		_, err := db.ExecContext(ctx, `INSERT INTO effort_entries (
			id,project_id,backlog_item_id,member_id,work_date,activity,hours_hundredths,sprint_id,deleted_at
		) VALUES (?, ?, 'a', 'member', '2026-10-20', 'Desarrollo', ?, ?, ?)`, id, projectValue.ID, hours, sprintID, deleted)
		if err != nil {
			t.Fatal(err)
		}
	}
	insertEffort("active-first", 600, first, nil)
	insertEffort("active-second", 100, second, nil)
	insertEffort("active-no-sprint", 75, nil, nil)
	insertEffort("deleted", 900, first, "2026-10-21T00:00:00Z")

	handler := httpapi.WithMetrics(
		httpapi.NewHandler(projectService, nil, sprintService),
		appmetrics.NewService(sqlite.NewMetricsRepository(db)),
	)
	response := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/metrics/hours", "")
	if response.Code != http.StatusOK {
		t.Fatalf("project metrics status = %d, body = %s", response.Code, response.Body)
	}
	var projectResult struct {
		Metrics domain.HoursMetrics `json:"metrics"`
		Sprints []struct {
			ID      int64               `json:"id"`
			Metrics domain.HoursMetrics `json:"metrics"`
		} `json:"sprints"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &projectResult); err != nil {
		t.Fatal(err)
	}
	if len(projectResult.Sprints) != 2 || !equalHourResult(projectResult.Metrics, 700, 775, 75) {
		t.Fatalf("project metrics = %+v", projectResult)
	}
	if projectResult.Sprints[0].ID != first || !equalHourResult(projectResult.Sprints[0].Metrics, 700, 600, -100) {
		t.Fatalf("first Sprint metrics = %+v", projectResult.Sprints[0])
	}
	if projectResult.Sprints[1].ID != second || !equalHourResult(projectResult.Sprints[1].Metrics, 400, 100, -300) {
		t.Fatalf("second Sprint metrics = %+v", projectResult.Sprints[1])
	}

	sprintResponse := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/sprints/"+itoa(first)+"/metrics/hours", "")
	var sprintResult domain.HoursMetrics
	if sprintResponse.Code != http.StatusOK || json.Unmarshal(sprintResponse.Body.Bytes(), &sprintResult) != nil || !equalHourResult(sprintResult, 700, 600, -100) {
		t.Fatalf("Sprint metrics = %d %s", sprintResponse.Code, sprintResponse.Body)
	}
	if got := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/sprints/nope/metrics/hours", "").Code; got != http.StatusBadRequest {
		t.Fatalf("invalid Sprint ID status = %d, want 400", got)
	}
}

func TestUS21MissingEstimateIsUnavailableAndDeletedEffortIsExcluded(t *testing.T) {
	ctx := context.Background()
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "Datos incompletos")
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), nil)
	value, err := sprintService.Create(ctx, projectValue.ID, "Sin estimación")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO members(id,project_id,name) VALUES('member',?,'Ana')`, projectValue.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO backlog_items(id,project_id,title,description,priority,status,acceptance_criteria)
		VALUES('a',?,'Historia','Descripción','Media','Pendiente','[]')`, projectValue.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO sprint_backlog_items(sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria)
		VALUES(?,'a','Historia','Descripción','Media','Pendiente','[]')`, value.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO effort_entries(id,project_id,backlog_item_id,member_id,work_date,activity,hours_hundredths,sprint_id,deleted_at)
		VALUES('active',?,'a','member','2026-10-20','Desarrollo',200,?,NULL),
		('deleted',?,'a','member','2026-10-20','Desarrollo',400,?,'2026-10-21T00:00:00Z')`, projectValue.ID, value.ID, projectValue.ID, value.ID); err != nil {
		t.Fatal(err)
	}
	handler := httpapi.WithMetrics(httpapi.NewHandler(projectService, nil, sprintService), appmetrics.NewService(sqlite.NewMetricsRepository(db)))
	response := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectValue.ID+"/sprints/"+itoa(value.ID)+"/metrics/hours", "")
	var result domain.HoursMetrics
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &result) != nil {
		t.Fatalf("metrics = %d %s", response.Code, response.Body)
	}
	if result.EstimatedHours != nil || result.Deviation != nil || result.ActualHours == nil || *result.ActualHours != 200 {
		t.Fatalf("metrics with missing estimate = %+v", result)
	}
}

func equalHourResult(value domain.HoursMetrics, estimated, actual, deviation int64) bool {
	return value.EstimatedHours != nil && int64(*value.EstimatedHours) == estimated &&
		value.ActualHours != nil && int64(*value.ActualHours) == actual &&
		value.Deviation != nil && int64(*value.Deviation) == deviation
}
