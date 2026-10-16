package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestBacklogHTTPCreateUpdateAndConsult(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	handler, projectID := newBacklogHTTPHandler(t, db)
	provisional := 5
	createBody := map[string]any{
		"title": "Login", "description": "A member can sign in", "priority": string(domain.PriorityHigh),
		"provisional_story_points": provisional, "acceptance_criteria": []string{"Valid credentials grant access"},
	}
	createdResponse := performJSONRequest(handler, http.MethodPost, "/api/projects/"+projectID+"/backlog", createBody)
	if createdResponse.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, body = %s; want %d", createdResponse.Code, createdResponse.Body, http.StatusCreated)
	}
	var created domain.Item
	if err := json.NewDecoder(createdResponse.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" || created.ProjectID != projectID || created.Status != domain.Pending {
		t.Errorf("created item = %+v, want server ID, project and Pending state", created)
	}
	if created.ProvisionalStoryPoints == nil || *created.ProvisionalStoryPoints != provisional || created.StoryPoints != nil {
		t.Errorf("created points = (%v, %v), want provisional %d and no definitive value", created.ProvisionalStoryPoints, created.StoryPoints, provisional)
	}

	updateBody := map[string]any{
		"title": "Updated login", "description": "Updated description", "priority": string(domain.PriorityMedium),
		"provisional_story_points": provisional, "acceptance_criteria": []string{"Updated criterion"}, "status": string(domain.Completed),
		"story_points": 99, "id": "client-id", "project_id": "other-project",
	}
	updatedResponse := performJSONRequest(handler, http.MethodPut, "/api/projects/"+projectID+"/backlog/"+created.ID, updateBody)
	if updatedResponse.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, body = %s; want %d", updatedResponse.Code, updatedResponse.Body, http.StatusOK)
	}
	var updated domain.Item
	if err := json.NewDecoder(updatedResponse.Body).Decode(&updated); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updated.ID != created.ID || updated.ProjectID != projectID || updated.Status != domain.Pending || updated.StoryPoints != nil {
		t.Errorf("PUT changed controlled fields: %+v", updated)
	}
	if updated.Title != "Updated login" || updated.ProvisionalStoryPoints == nil || *updated.ProvisionalStoryPoints != provisional {
		t.Errorf("updated item = %+v, want editable fields changed and provisional points preserved", updated)
	}

	listResponse := performJSONRequest(handler, http.MethodGet, "/api/projects/"+projectID+"/backlog", nil)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("GET list status = %d, body = %s; want %d", listResponse.Code, listResponse.Body, http.StatusOK)
	}
	var items []domain.Item
	if err := json.NewDecoder(listResponse.Body).Decode(&items); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(items) != 1 || !reflect.DeepEqual(items[0], updated) {
		t.Errorf("GET list = %+v, want updated item %+v", items, updated)
	}
}

func TestBacklogHTTPEmptyAndMissingProject(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	handler, projectID := newBacklogHTTPHandler(t, db)
	empty := performJSONRequest(handler, http.MethodGet, "/api/projects/"+projectID+"/backlog", nil)
	if empty.Code != http.StatusOK || empty.Body.String() != "[]\n" {
		t.Errorf("empty GET = (%d, %s), want 200 and []", empty.Code, empty.Body)
	}
	missing := performJSONRequest(handler, http.MethodGet, "/api/projects/missing/backlog", nil)
	if missing.Code != http.StatusNotFound {
		t.Errorf("missing GET status = %d, want %d", missing.Code, http.StatusNotFound)
	}
}

func TestBacklogHTTPValidationAndScopedUpdate(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	handler, projectID := newBacklogHTTPHandler(t, db)
	invalid := performJSONRequest(handler, http.MethodPost, "/api/projects/"+projectID+"/backlog", map[string]any{
		"title": " ", "description": "Description", "priority": "Alta", "acceptance_criteria": []string{"criterion"},
	})
	if invalid.Code != http.StatusBadRequest {
		t.Errorf("invalid POST status = %d, want %d", invalid.Code, http.StatusBadRequest)
	}
	created := performJSONRequest(handler, http.MethodPost, "/api/projects/"+projectID+"/backlog", map[string]any{
		"title": "Login", "description": "Description", "priority": "Alta", "acceptance_criteria": []string{"criterion"},
	})
	var item domain.Item
	if err := json.NewDecoder(created.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	otherProject, err := projects.NewService(sqlite.NewProjectRepository(db)).Create(context.Background(), "Other")
	if err != nil {
		t.Fatal(err)
	}
	otherProjectUpdate := performJSONRequest(handler, http.MethodPut, "/api/projects/"+otherProject.ID+"/backlog/"+item.ID, map[string]any{
		"title": "Wrong project", "description": "Description", "priority": "Alta", "acceptance_criteria": []string{"criterion"},
	})
	if otherProjectUpdate.Code != http.StatusNotFound {
		t.Errorf("cross-project PUT status = %d, want %d", otherProjectUpdate.Code, http.StatusNotFound)
	}
}

func TestBacklogHTTPReportsStorageFailure(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	_, projectID := newBacklogHTTPHandler(t, db)
	handler := httpapi.NewHandler(
		projects.NewService(sqlite.NewProjectRepository(db)),
		appbacklog.NewService(sqlite.NewBacklogRepository(db)),
		nil,
	)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	response := performJSONRequest(handler, http.MethodGet, "/api/projects/"+projectID+"/backlog", nil)
	if response.Code != http.StatusInternalServerError {
		t.Errorf("storage failure status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}

func newBacklogHTTPHandler(t *testing.T, db *sql.DB) (http.Handler, string) {
	t.Helper()
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	created, err := projectService.Create(context.Background(), "Carva")
	if err != nil {
		t.Fatal(err)
	}
	backlogService := appbacklog.NewService(sqlite.NewBacklogRepository(db))
	return httpapi.NewHandler(projectService, backlogService, nil), created.ID
}

func performJSONRequest(handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&requestBody).Encode(body)
	}
	request := httptest.NewRequest(method, path, &requestBody)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
