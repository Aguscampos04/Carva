package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/members"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	sqlite "github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestHUAssigneePersistsAndRejectsForeignMember(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	ctx := context.Background()
	ps := projects.NewService(sqlite.NewProjectRepository(db))
	p, _ := ps.Create(ctx, "A")
	other, _ := ps.Create(ctx, "B")
	ms := members.NewService(sqlite.NewMemberRepository(db))
	member, err := ms.Create(ctx, p.ID, "Ana")
	if err != nil {
		t.Fatal(err)
	}
	foreign, err := ms.Create(ctx, other.ID, "Ana")
	if err != nil {
		t.Fatal(err)
	}
	repo := sqlite.NewBacklogRepository(db)
	estimate := hours.Value(100)
	story, err := backlog.NewService(repo).Create(ctx, p.ID, domain.Input{Title: "HU", Description: "Detail", Priority: domain.PriorityHigh, EstimatedHours: &estimate, AcceptanceCriteria: []string{"OK"}})
	if err != nil {
		t.Fatal(err)
	}
	api := httpapi.WithAssignees(httpapi.NewHandler(ps, nil, nil), backlog.NewAssigneeService(sqlite.NewAssigneeRepository(db)))
	url := "/api/projects/" + p.ID + "/backlog/" + story.ID + "/assignee"
	call := func(method, path, body string) *httptest.ResponseRecorder {
		r := httptest.NewRecorder()
		api.ServeHTTP(r, httptest.NewRequest(method, path, strings.NewReader(body)))
		return r
	}
	if r := call("GET", url, ""); r.Code != 200 || strings.TrimSpace(r.Body.String()) != "{\"member_id\":null}" {
		t.Fatalf("initial: %d %s", r.Code, r.Body)
	}
	for _, body := range []string{"{}", "{\"member_id\":\"bad\"}", "{\"member_id\":\"" + foreign.ID + "\"}"} {
		if r := call("PUT", url, body); r.Code != 400 {
			t.Fatalf("invalid assignment: %d %s", r.Code, r.Body)
		}
	}
	if r := call("PUT", url, "{\"member_id\":\""+member.ID+"\"}"); r.Code != 200 {
		t.Fatalf("save: %d %s", r.Code, r.Body)
	}
	if r := call("GET", url, ""); !strings.Contains(r.Body.String(), member.ID) {
		t.Fatalf("not persisted: %s", r.Body)
	}
	ss := sprints.NewService(sqlite.NewSprintRepository(db), repo)
	first, err := ss.Create(ctx, p.ID, "Primero")
	if err != nil {
		t.Fatal(err)
	}
	next, err := ss.Create(ctx, p.ID, "Siguiente")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ss.AssignBacklogItem(ctx, p.ID, first.ID, story.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := ss.Close(ctx, p.ID, first.ID); err != nil {
		t.Fatal(err)
	}
	planned, err := ss.ListPlannedBacklogItems(ctx, p.ID, next.ID)
	if err != nil || len(planned) != 1 {
		t.Fatalf("transfer: %+v %v", planned, err)
	}
	if r := call("GET", url, ""); !strings.Contains(r.Body.String(), member.ID) {
		t.Fatal("responsible member lost after Sprint transfer")
	}
	foreignURL := strings.Replace(url, p.ID, other.ID, 1)
	if r := call("GET", foreignURL, ""); r.Code != 404 {
		t.Fatalf("foreign story: %d", r.Code)
	}
	if r := call("PUT", url, "{\"member_id\":null}"); r.Code != 200 {
		t.Fatalf("unassign: %d %s", r.Code, r.Body)
	}
	if r := call("GET", url, ""); !strings.Contains(r.Body.String(), "null") {
		t.Fatal("unassign not persisted")
	}
	_ = db.Close()
	if r := call(http.MethodGet, url, ""); r.Code != 500 || strings.Contains(r.Body.String(), "sql:") {
		t.Fatalf("storage failure: %d %s", r.Code, r.Body)
	}
}
