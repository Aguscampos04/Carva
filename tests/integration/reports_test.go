package integration_test

import (
	"bytes"
	"encoding/json"
	appmetrics "github.com/Aguscampos04/Carva/internal/application/metrics"
	app "github.com/Aguscampos04/Carva/internal/application/reports"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
	"net/http"
	"strings"
	"testing"
)

func TestUS24ReportsMatchSourcesAndPreserveScope(t *testing.T) {
	f := newDefectUpdateFixture(t)
	service := app.NewService(sqlite.NewReportSource(f.db), appmetrics.NewService(sqlite.NewMetricsRepository(f.db)))
	handler := httpapi.WithReports(f.handler, service)
	base := "/api/projects/p/reports/"
	get := func(path string) app.Report {
		t.Helper()
		response := performJSON(t, handler, "GET", path, "")
		if response.Code != 200 {
			t.Fatalf("%s: %d %s", path, response.Code, response.Body)
		}
		var report app.Report
		if err := json.Unmarshal(response.Body.Bytes(), &report); err != nil {
			t.Fatal(err)
		}
		return report
	}
	project := get(base + "project")
	if len(project.Stories) != 2 || len(project.Sprints) != 2 || len(project.Defects) != 1 || project.Scope != "project" {
		t.Fatalf("report %+v", project)
	}
	metrics := performJSON(t, httpapi.WithMetrics(f.handler, appmetrics.NewService(sqlite.NewMetricsRepository(f.db))), "GET", "/api/projects/p/metrics/defects", "")
	var source struct {
		Metrics json.RawMessage `json:"metrics"`
	}
	json.Unmarshal(metrics.Body.Bytes(), &source)
	actual, _ := json.Marshal(project.Metrics.Defects)
	if !bytes.Equal(actual, source.Metrics) {
		t.Fatalf("source %s != %s", source.Metrics, actual)
	}
	patch := performJSON(t, handler, "PATCH", f.url, `{"status":"Resuelto","detection_sprint_id":1,"resolution_sprint_id":2}`)
	expectCode(t, patch, 200)
	first, second := get(base+"sprints/1"), get(base+"sprints/2")
	if len(first.Defects) != 1 || len(second.Defects) != 1 || len(first.Stories) != 0 || len(first.Sprints) != 1 {
		t.Fatal("scope or split attribution")
	}
	if *first.Metrics.Defects.DetectedDefects != 1 || *first.Metrics.Defects.ResolvedDefects != 0 || *second.Metrics.Defects.ResolvedDefects != 1 {
		t.Fatal("incorrect source metrics")
	}
	for _, path := range []string{base + "project.pdf", base + "sprints/1.pdf"} {
		response := performJSON(t, handler, http.MethodGet, path, "")
		expectCode(t, response, 200)
		if response.Header().Get("Content-Type") != "application/pdf" || !strings.Contains(response.Header().Get("Content-Disposition"), "attachment") || !bytes.HasPrefix(response.Body.Bytes(), []byte("%PDF-")) || !bytes.Contains(response.Body.Bytes(), []byte("%%EOF")) {
			t.Fatalf("invalid PDF %s", response.Body)
		}
	}
	// Nothing in report/export changes the stored defect.
	if got := listedDefect(t, f.handler); got["status"] != "Resuelto" || got["detection_sprint_id"] != float64(1) {
		t.Fatal(got)
	}
}

func TestUS24MissingInvalidAndUnavailableSources(t *testing.T) {
	f := newDefectUpdateFixture(t)
	handler := httpapi.WithReports(f.handler, app.NewService(sqlite.NewReportSource(f.db), appmetrics.NewService(sqlite.NewMetricsRepository(f.db))))
	for _, tt := range []struct {
		path string
		code int
	}{{"/api/projects/missing/reports/project", 404}, {"/api/projects/p/reports/sprints/3", 404}, {"/api/projects/p/reports/sprints/999.pdf", 404}, {"/api/projects/p/reports/sprints/0", 400}, {"/api/projects/p/reports/sprints/nope.pdf", 400}} {
		expectCode(t, performJSON(t, handler, "GET", tt.path, ""), tt.code)
	}
	f.db.Close()
	response := performJSON(t, handler, "GET", "/api/projects/p/reports/project.pdf", "")
	expectCode(t, response, 500)
	if bytes.HasPrefix(response.Body.Bytes(), []byte("%PDF")) {
		t.Fatal("error reported as PDF")
	}
}

func TestUS24HistoricalPlanningAndEffortAreNotRebuiltOrDoubleCounted(t *testing.T) {
	f := newDefectUpdateFixture(t)
	for _, query := range []string{
		`UPDATE backlog_items SET title='Changed current title',status='Completada',story_points=5,estimated_hours_hundredths=500 WHERE id='a'`,
		`INSERT INTO sprint_backlog_items(sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria,story_points,completion_verified,estimated_hours_hundredths) VALUES(1,'a','Original planning','Original description','Alta','Pendiente','["Original criterion"]',5,0,500),(2,'a','Next planning','Next description','Alta','Pendiente','["Next criterion"]',5,1,500)`,
		`INSERT INTO effort_entries(id,project_id,backlog_item_id,member_id,work_date,activity,hours_hundredths,sprint_id,deleted_at) VALUES('active','p','a','m','2026-10-20','Active work',375,1,NULL),('unassigned','p','a','m','2026-10-20','Without Sprint',125,NULL,NULL),('removed','p','a','m','2026-10-20','Deleted work',200,1,'2026-10-21')`,
	} {
		if _, err := f.db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	handler := httpapi.WithReports(f.handler, app.NewService(sqlite.NewReportSource(f.db), appmetrics.NewService(sqlite.NewMetricsRepository(f.db))))
	get := func(path string) app.Report {
		t.Helper()
		r := performJSON(t, handler, "GET", path, "")
		expectCode(t, r, 200)
		var result app.Report
		if err := json.Unmarshal(r.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	project := get("/api/projects/p/reports/project")
	if len(project.Effort) != 2 || *project.Metrics.Hours.ActualHours != 500 {
		t.Fatal("duplicate, missing or deleted effort", project.Effort, project.Metrics.Hours)
	}
	first := get("/api/projects/p/reports/sprints/1")
	if len(first.Effort) != 1 || *first.Metrics.Hours.ActualHours != 375 {
		t.Fatal("incorrect scoped effort")
	}
	planned := first.Sprints[0].Stories[0]
	if planned.Title != "Original planning" || planned.AcceptanceCriteriaVerified || planned.AcceptanceCriteria[0] != "Original criterion" || *first.Metrics.Points.CompletedStoryPoints != 0 {
		t.Fatal("history overwritten", planned)
	}
	next := get("/api/projects/p/reports/sprints/2")
	if !next.Sprints[0].Stories[0].AcceptanceCriteriaVerified || *next.Metrics.Points.CompletedStoryPoints != 5 || len(next.Effort) != 0 {
		t.Fatal("destination or scope incorrect")
	}
	for _, path := range []string{"/api/projects/p/reports/1", "/api/projects/p/reports/unknown.pdf"} {
		expectCode(t, performJSON(t, handler, "GET", path, ""), 404)
	}
	expectCode(t, performJSON(t, handler, "GET", "/api/projects/p/reports/sprints/project", ""), 400)
	// Repeat consultation to verify persisted original data remains unchanged.
	again := get("/api/projects/p/reports/sprints/1")
	if again.Sprints[0].Stories[0].Title != planned.Title {
		t.Fatal("consultation changed planning")
	}
}
