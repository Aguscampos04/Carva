package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/members"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestUS56CompletionRequiresEffortAndRetriesAtomically(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	ctx := context.Background()
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "US-56")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := appbacklog.NewService(backlogRepository)
	item, err := backlogService.Create(ctx, projectValue.ID, backlog.Input{
		EstimatedHours: hours.Ptr(800), Title: "Historia", Description: "Descripción",
		Priority: backlog.PriorityHigh, AcceptanceCriteria: []string{"Criterio"},
	})
	if err != nil {
		t.Fatal(err)
	}
	memberService := members.NewService(sqlite.NewMemberRepository(db))
	registered, err := memberService.Create(ctx, projectValue.ID, "Ana")
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), backlogRepository)
	sprintValue, err := sprintService.Create(ctx, projectValue.ID, "Completar historia")
	if err != nil {
		t.Fatal(err)
	}
	handler := httpapi.NewHandler(projectService, backlogService, sprintService, memberService)
	assignmentPath := "/api/projects/" + projectValue.ID + "/sprints/" + strconv.FormatInt(sprintValue.ID, 10) + "/backlog-items"
	if response := performJSONRequest(handler, http.MethodPost, assignmentPath, map[string]any{"backlog_item_id": item.ID}); response.Code != http.StatusCreated {
		t.Fatalf("assign story = %d %s", response.Code, response.Body)
	}
	completionPath := assignmentPath + "/" + item.ID + "/completion"
	missing := performJSONRequest(handler, http.MethodPost, completionPath, map[string]any{"acceptance_criteria_verified": true})
	if missing.Code != http.StatusConflict {
		t.Fatalf("completion without effort = %d %s; want conflict", missing.Code, missing.Body)
	}
	var missingBody map[string]any
	if err := json.Unmarshal(missing.Body.Bytes(), &missingBody); err != nil || missingBody["code"] != "effort_required" {
		t.Fatalf("missing-effort response = %s; want effort_required", missing.Body)
	}
	assertUS56StoryPending(t, db, projectValue.ID, item.ID)

	workDate := us56Today(t)
	effort := map[string]any{"member_id": registered.ID, "work_date": workDate, "activity": "Implementación", "hours": 1.25}
	request := map[string]any{"acceptance_criteria_verified": true, "effort": effort}
	if _, err := db.Exec(`CREATE TRIGGER reject_us56_completion BEFORE UPDATE OF completion_verified ON sprint_backlog_items BEGIN SELECT RAISE(ABORT, 'forced completion failure'); END`); err != nil {
		t.Fatal(err)
	}
	failed := performJSONRequest(handler, http.MethodPost, completionPath, request)
	if failed.Code != http.StatusInternalServerError {
		t.Fatalf("forced completion error = %d %s; want 500", failed.Code, failed.Body)
	}
	assertUS56StoryPending(t, db, projectValue.ID, item.ID)
	assertUS56EffortCount(t, db, projectValue.ID, item.ID, 0)
	if _, err := db.Exec(`DROP TRIGGER reject_us56_completion`); err != nil {
		t.Fatal(err)
	}

	completed := performJSONRequest(handler, http.MethodPost, completionPath, request)
	if completed.Code != http.StatusOK {
		t.Fatalf("valid completion = %d %s; want 200", completed.Code, completed.Body)
	}
	retried := performJSONRequest(handler, http.MethodPost, completionPath, request)
	if retried.Code != http.StatusOK {
		t.Fatalf("repeated completion = %d %s; want idempotent success", retried.Code, retried.Body)
	}
	assertUS56EffortCount(t, db, projectValue.ID, item.ID, 1)
	assertUS56StoryCompleted(t, db, projectValue.ID, item.ID)
}

func TestUS56ExistingActiveEffortSatisfiesCompletionWithoutExtraEntry(t *testing.T) {
	f := newSprintClosureFixture(t)
	memberService := members.NewService(sqlite.NewMemberRepository(f.db))
	registered, err := memberService.Create(context.Background(), f.projectID, "Ana")
	if err != nil {
		t.Fatal(err)
	}
	item := f.createStory(t, "Existing effort")
	f.assign(t, f.current, item.ID)
	today := us56Today(t)
	if _, err := f.db.Exec(`INSERT INTO effort_entries(id,project_id,backlog_item_id,member_id,work_date,activity,hours_hundredths,sprint_id) VALUES('effort-old',?,?,?,?,?,?,?)`, f.projectID, item.ID, registered.ID, today, "Trabajo ya cargado", 125, f.current.ID); err != nil {
		t.Fatal(err)
	}
	path := sprintPath(f.projectID, f.current.ID) + "/backlog-items/" + item.ID + "/completion"
	response := performJSONRequest(f.handler(), http.MethodPost, path, map[string]any{"acceptance_criteria_verified": true})
	if response.Code != http.StatusOK {
		t.Fatalf("completion with existing effort = %d %s; want 200", response.Code, response.Body)
	}
	assertUS56EffortCount(t, f.db, f.projectID, item.ID, 1)
	assertUS56StoryCompleted(t, f.db, f.projectID, item.ID)
}

func assertUS56EffortCount(t *testing.T, db interface{ QueryRow(string, ...any) *sql.Row }, projectID, itemID string, want int) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM effort_entries WHERE project_id=? AND backlog_item_id=? AND deleted_at IS NULL`, projectID, itemID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Fatalf("active effort rows = %d, want %d", count, want)
	}
}

func assertUS56StoryPending(t *testing.T, db interface{ QueryRow(string, ...any) *sql.Row }, projectID, itemID string) {
	t.Helper()
	var status string
	if err := db.QueryRow(`SELECT status FROM backlog_items WHERE project_id=? AND id=?`, projectID, itemID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(backlog.Pending) {
		t.Fatalf("Backlog status = %q, want %q", status, backlog.Pending)
	}
}

func assertUS56StoryCompleted(t *testing.T, db interface{ QueryRow(string, ...any) *sql.Row }, projectID, itemID string) {
	t.Helper()
	var status string
	if err := db.QueryRow(`SELECT status FROM backlog_items WHERE project_id=? AND id=?`, projectID, itemID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(backlog.Completed) {
		t.Fatalf("Backlog status = %q, want %q", status, backlog.Completed)
	}
}

func us56Today(t *testing.T) string {
	t.Helper()
	location, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		t.Fatal(err)
	}
	return time.Now().In(location).Format("2006-01-02")
}
