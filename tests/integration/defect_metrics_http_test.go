package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	appmetrics "github.com/Aguscampos04/Carva/internal/application/metrics"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/domain/defect"
	domain "github.com/Aguscampos04/Carva/internal/domain/metrics"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestUS22DefectMetricsHTTPSeparatesDetectionResolutionAndProjectTotals(t *testing.T) {
	fixture := newDefectMetricsHTTPFixture(t)
	fixture.insert("split", defect.StatusResolved, fixture.first, fixture.second)
	fixture.insert("open", defect.StatusOpen, fixture.second, nil)

	response := performJSON(t, fixture.handler, http.MethodGet, "/api/projects/"+fixture.projectID+"/metrics/defects", "")
	var result domain.ProjectDefectMetrics
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &result) != nil {
		t.Fatalf("project defect metrics = %d %s", response.Code, response.Body)
	}
	if len(result.Sprints) != 2 || *result.Metrics.DetectedDefects != 2 || *result.Metrics.ResolvedDefects != 1 {
		t.Fatalf("project defect metrics = %+v", result)
	}
	if *result.Sprints[0].Metrics.DetectedDefects != 1 || *result.Sprints[0].Metrics.ResolvedDefects != 0 ||
		*result.Sprints[1].Metrics.DetectedDefects != 1 || *result.Sprints[1].Metrics.ResolvedDefects != 1 {
		t.Fatalf("Sprint defect metrics = %+v", result.Sprints)
	}

	response = performJSON(t, fixture.handler, http.MethodGet,
		"/api/projects/"+fixture.projectID+"/sprints/"+itoa(fixture.second)+"/metrics/defects", "")
	var sprintResult domain.DefectMetrics
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &sprintResult) != nil ||
		*sprintResult.DetectedDefects != 1 || *sprintResult.ResolvedDefects != 1 {
		t.Fatalf("Sprint endpoint metrics = %d %s", response.Code, response.Body)
	}
}

func TestUS22DefectMetricsHTTPReportsUnavailableAttribution(t *testing.T) {
	fixture := newDefectMetricsHTTPFixture(t)
	fixture.insert("unassigned-detection", defect.StatusOpen, nil, nil)
	fixture.insert("unassigned-resolution", defect.StatusClosed, fixture.first, nil)

	response := performJSON(t, fixture.handler, http.MethodGet, "/api/projects/"+fixture.projectID+"/metrics/defects", "")
	var result domain.ProjectDefectMetrics
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &result) != nil {
		t.Fatalf("project defect metrics = %d %s", response.Code, response.Body)
	}
	if *result.Metrics.DetectedDefects != 2 || *result.Metrics.ResolvedDefects != 1 {
		t.Fatalf("project totals = %+v, want 2 detected and 1 resolved", result.Metrics)
	}
	if result.Sprints[0].Metrics.DetectedDefects != nil || result.Sprints[1].Metrics.DetectedDefects != nil ||
		result.Sprints[0].Metrics.ResolvedDefects != nil || result.Sprints[1].Metrics.ResolvedDefects != nil {
		t.Fatalf("Sprint metrics with missing attribution = %+v, want null values", result.Sprints)
	}
}

func TestUS22DefectMetricsHTTPValidatesSprintAndProject(t *testing.T) {
	fixture := newDefectMetricsHTTPFixture(t)
	base := "/api/projects/" + fixture.projectID
	for _, test := range []struct {
		path string
		want int
	}{
		{path: "/api/projects/missing/metrics/defects", want: http.StatusNotFound},
		{path: base + "/sprints/999/metrics/defects", want: http.StatusNotFound},
		{path: base + "/sprints/nope/metrics/defects", want: http.StatusBadRequest},
		{path: base + "/sprints/0/metrics/defects", want: http.StatusBadRequest},
	} {
		if got := performJSON(t, fixture.handler, http.MethodGet, test.path, "").Code; got != test.want {
			t.Errorf("GET %s status = %d, want %d", test.path, got, test.want)
		}
	}
}

type defectMetricsHTTPFixture struct {
	handler   http.Handler
	projectID string
	first     int64
	second    int64
	insert    func(string, defect.Status, any, any)
}

func newDefectMetricsHTTPFixture(t *testing.T) defectMetricsHTTPFixture {
	t.Helper()
	ctx := context.Background()
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "Defect metrics HTTP")
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), nil)
	first, err := sprintService.Create(ctx, projectValue.ID, "First")
	if err != nil {
		t.Fatal(err)
	}
	second, err := sprintService.Create(ctx, projectValue.ID, "Second")
	if err != nil {
		t.Fatal(err)
	}
	base := httpapi.NewHandler(projectService, nil, sprintService)
	handler := httpapi.WithMetrics(base, appmetrics.NewService(sqlite.NewMetricsRepository(db)))
	insert := func(id string, status defect.Status, detection, resolution any) {
		t.Helper()
		storyID := "story-" + id
		if _, err := db.ExecContext(ctx, `INSERT INTO backlog_items (
			id,project_id,title,description,priority,status,acceptance_criteria
		) VALUES (?,?,'Story','Description','Media','Pendiente','[]')`, storyID, projectValue.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := db.ExecContext(ctx, `INSERT INTO defects (
			id,project_id,backlog_item_id,description,severity,status,detection_sprint_id,resolution_sprint_id
		) VALUES (?,?,?,'Defect','Media',?,?,?)`, id, projectValue.ID, storyID, status, detection, resolution); err != nil {
			t.Fatal(err)
		}
	}
	return defectMetricsHTTPFixture{handler: handler, projectID: projectValue.ID, first: first.ID, second: second.ID, insert: insert}
}
