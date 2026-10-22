package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestUS10CloseSprintHTTPReturnsPersistedClosedState(t *testing.T) {
	fixture := newSprintClosureFixture(t)
	handler := fixture.handler()
	path := sprintPath(fixture.projectID, fixture.current.ID) + "/close"

	closed := performJSONRequest(handler, http.MethodPost, path, nil)
	if closed.Code != http.StatusOK {
		t.Fatalf("close Sprint status = %d, body = %s; want %d", closed.Code, closed.Body, http.StatusOK)
	}
	var result struct {
		ID     int64 `json:"id"`
		Closed bool  `json:"closed"`
	}
	if err := json.NewDecoder(closed.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.ID != fixture.current.ID || !result.Closed {
		t.Fatalf("close response = %+v, want closed Sprint %d", result, fixture.current.ID)
	}

	consulted := performJSONRequest(handler, http.MethodGet, sprintPath(fixture.projectID, fixture.current.ID), nil)
	var persisted struct {
		Closed bool `json:"closed"`
	}
	if consulted.Code != http.StatusOK || json.NewDecoder(consulted.Body).Decode(&persisted) != nil || !persisted.Closed {
		t.Fatalf("consult closed Sprint = %d %s, want persisted closed=true", consulted.Code, consulted.Body)
	}
}

func TestUS10CloseSprintHTTPBlocksWithoutNextAndClosedMutations(t *testing.T) {
	fixture := newSprintClosureFixture(t)
	handler := fixture.handler()
	item := fixture.createStory(t, "Incomplete")
	fixture.assign(t, fixture.current, item.ID)
	if _, err := fixture.db.Exec(`DELETE FROM sprints WHERE id = ?`, fixture.next.ID); err != nil {
		t.Fatal(err)
	}
	closePath := sprintPath(fixture.projectID, fixture.current.ID) + "/close"
	blocked := performJSONRequest(handler, http.MethodPost, closePath, nil)
	if blocked.Code != http.StatusConflict {
		t.Fatalf("close without next Sprint status = %d, body = %s; want %d", blocked.Code, blocked.Body, http.StatusConflict)
	}
	current, err := fixture.repo.Get(context.Background(), fixture.projectID, fixture.current.ID)
	if err != nil || current.Closed {
		t.Fatalf("Sprint after blocked close = %+v, %v; want open", current, err)
	}

	fixture.next, err = fixture.service.Create(context.Background(), fixture.projectID, "Next Sprint")
	if err != nil {
		t.Fatal(err)
	}
	closed := performJSONRequest(handler, http.MethodPost, closePath, nil)
	if closed.Code != http.StatusOK {
		t.Fatalf("close with next Sprint status = %d, body = %s", closed.Code, closed.Body)
	}
	newItem := fixture.createStory(t, "New story")
	assignPath := sprintPath(fixture.projectID, fixture.current.ID) + "/backlog-items"
	assigned := performJSONRequest(handler, http.MethodPost, assignPath, map[string]any{"backlog_item_id": newItem.ID})
	if assigned.Code != http.StatusConflict {
		t.Fatalf("assign to closed Sprint status = %d; want %d", assigned.Code, http.StatusConflict)
	}
	completionPath := assignPath + "/" + item.ID + "/completion"
	completed := performJSONRequest(handler, http.MethodPost, completionPath, map[string]any{"acceptance_criteria_verified": true})
	if completed.Code != http.StatusConflict {
		t.Fatalf("complete story on closed Sprint status = %d; want %d", completed.Code, http.StatusConflict)
	}
	repeated := performJSONRequest(handler, http.MethodPost, closePath, nil)
	if repeated.Code != http.StatusConflict {
		t.Fatalf("repeated close status = %d; want %d", repeated.Code, http.StatusConflict)
	}
}

func (fixture *sprintClosureFixture) handler() http.Handler {
	projectService := projects.NewService(sqlite.NewProjectRepository(fixture.db))
	backlogRepository := sqlite.NewBacklogRepository(fixture.db)
	sprintService := sprints.NewService(fixture.repo, backlogRepository)
	return httpapi.NewHandler(projectService, nil, sprintService)
}

func sprintPath(projectID string, sprintID int64) string {
	return "/api/projects/" + projectID + "/sprints/" + strconv.FormatInt(sprintID, 10)
}
