package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestCreateProjectAndFirstSprintThroughAPI(t *testing.T) {
	handler := newTestHandler(t)

	projectResponse := performJSON(t, handler, http.MethodPost, "/api/projects", `{"name":"Carva Demo"}`)
	if projectResponse.Code != http.StatusCreated {
		t.Fatalf("create project status = %d, body = %s", projectResponse.Code, projectResponse.Body)
	}
	var createdProject struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(projectResponse.Body.Bytes(), &createdProject); err != nil {
		t.Fatalf("decode project response: %v", err)
	}
	if createdProject.ID == "" || createdProject.Name != "Carva Demo" || createdProject.Status != "Planificado" {
		t.Fatalf("created project = %+v", createdProject)
	}

	sprintPath := "/api/projects/" + createdProject.ID + "/sprints"
	sprintResponse := performJSON(t, handler, http.MethodPost, sprintPath, `{"goal":"Entregar el MVP"}`)
	if sprintResponse.Code != http.StatusCreated {
		t.Fatalf("create sprint status = %d, body = %s", sprintResponse.Code, sprintResponse.Body)
	}
	var createdSprint struct {
		ID        int64  `json:"id"`
		ProjectID string `json:"projectID"`
		Goal      string `json:"goal"`
	}
	if err := json.Unmarshal(sprintResponse.Body.Bytes(), &createdSprint); err != nil {
		t.Fatalf("decode sprint response: %v", err)
	}
	if createdSprint.ID < 1 || createdSprint.ProjectID != createdProject.ID || createdSprint.Goal != "Entregar el MVP" {
		t.Fatalf("created sprint = %+v", createdSprint)
	}

	projectGet := performJSON(t, handler, http.MethodGet, "/api/projects/"+createdProject.ID, "")
	if projectGet.Code != http.StatusOK {
		t.Fatalf("get project status = %d, body = %s", projectGet.Code, projectGet.Body)
	}
	var updatedProject struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(projectGet.Body.Bytes(), &updatedProject); err != nil {
		t.Fatalf("decode updated project response: %v", err)
	}
	if updatedProject.Status != "En progreso" {
		t.Fatalf("project status = %q, want En progreso", updatedProject.Status)
	}

	sprintGet := performJSON(t, handler, http.MethodGet,
		"/api/projects/"+createdProject.ID+"/sprints/"+itoa(createdSprint.ID), "")
	if sprintGet.Code != http.StatusOK {
		t.Fatalf("get sprint status = %d, body = %s", sprintGet.Code, sprintGet.Body)
	}
	var fetchedSprint struct {
		Goal string `json:"goal"`
	}
	if err := json.Unmarshal(sprintGet.Body.Bytes(), &fetchedSprint); err != nil {
		t.Fatalf("decode fetched sprint response: %v", err)
	}
	if fetchedSprint.Goal != "Entregar el MVP" {
		t.Fatalf("fetched sprint goal = %q", fetchedSprint.Goal)
	}
}

func TestInvalidSprintCreationDoesNotChangeProject(t *testing.T) {
	handler := newTestHandler(t)
	projectResponse := performJSON(t, handler, http.MethodPost, "/api/projects", `{"name":"Carva Demo"}`)
	var createdProject struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(projectResponse.Body.Bytes(), &createdProject); err != nil {
		t.Fatalf("decode project response: %v", err)
	}

	response := performJSON(t, handler, http.MethodPost, "/api/projects/"+createdProject.ID+"/sprints", `{"goal":"  "}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("blank goal status = %d, want 400; body = %s", response.Code, response.Body)
	}
	projectGet := performJSON(t, handler, http.MethodGet, "/api/projects/"+createdProject.ID, "")
	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(projectGet.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode project response: %v", err)
	}
	if result.Status != "Planificado" {
		t.Fatalf("project status after rejected sprint = %q, want Planificado", result.Status)
	}
}

func TestCreateSprintRejectsUnknownProject(t *testing.T) {
	handler := newTestHandler(t)
	response := performJSON(t, handler, http.MethodPost, "/api/projects/unknown-project/sprints", `{"goal":"Entregar el MVP"}`)
	if response.Code != http.StatusNotFound {
		t.Fatalf("unknown project status = %d, want 404; body = %s", response.Code, response.Body)
	}
}

func TestSecondSprintKeepsProjectInProgress(t *testing.T) {
	handler := newTestHandler(t)
	projectID := createProject(t, handler, "Carva Demo")
	for _, goal := range []string{"First goal", "Second goal"} {
		response := performJSON(t, handler, http.MethodPost, "/api/projects/"+projectID+"/sprints", `{"goal":"`+goal+`"}`)
		if response.Code != http.StatusCreated {
			t.Fatalf("create Sprint %q status = %d, body = %s", goal, response.Code, response.Body)
		}
	}
	projectResponse := performJSON(t, handler, http.MethodGet, "/api/projects/"+projectID, "")
	var value struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(projectResponse.Body.Bytes(), &value); err != nil {
		t.Fatalf("decode project response: %v", err)
	}
	if value.Status != "En progreso" {
		t.Fatalf("project status = %q, want En progreso", value.Status)
	}
}

func TestSprintCannotBeQueriedThroughAnotherProject(t *testing.T) {
	handler := newTestHandler(t)
	first := createProject(t, handler, "First")
	second := createProject(t, handler, "Second")
	response := performJSON(t, handler, http.MethodPost, "/api/projects/"+first+"/sprints", `{"goal":"First goal"}`)
	var createdSprint struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &createdSprint); err != nil {
		t.Fatalf("decode sprint response: %v", err)
	}
	get := performJSON(t, handler, http.MethodGet,
		"/api/projects/"+second+"/sprints/"+itoa(createdSprint.ID), "")
	if get.Code != http.StatusNotFound {
		t.Fatalf("cross-project sprint lookup status = %d, want 404", get.Code)
	}
}

func TestCreateProjectRejectsBlankName(t *testing.T) {
	handler := newTestHandler(t)
	response := performJSON(t, handler, http.MethodPost, "/api/projects", `{"name":"  "}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("blank project name status = %d, want 400", response.Code)
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	db, err := sqlite.OpenPath(context.Background(), filepath.Join(t.TempDir(), "carva.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := sqlite.Migrate(context.Background(), db); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	projectRepository := sqlite.NewProjectRepository(db)
	sprintRepository := sqlite.NewSprintRepository(db)
	return httpapi.NewHandler(projects.NewService(projectRepository), sprints.NewService(sprintRepository))
}

func createProject(t *testing.T, handler http.Handler, name string) string {
	t.Helper()
	response := performJSON(t, handler, http.MethodPost, "/api/projects", `{"name":"`+name+`"}`)
	if response.Code != http.StatusCreated {
		t.Fatalf("create project status = %d, body = %s", response.Code, response.Body)
	}
	var body struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode project response: %v", err)
	}
	return body.ID
}

func performJSON(t *testing.T, handler http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func itoa(value int64) string {
	return strconv.FormatInt(value, 10)
}
