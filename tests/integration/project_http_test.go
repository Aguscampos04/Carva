package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	carvasqlite "github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestProjectHTTPCreateReadAndRejectClientControlledIdentityOrStatus(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	handler := httpapi.NewHandler(projects.NewService(carvasqlite.NewProjectRepository(db)), nil, nil)

	request := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(
		`{"name":"Carva MVP","id":"client-id","status":"Finalizado"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	created := httptest.NewRecorder()
	handler.ServeHTTP(created, request)
	if created.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, body = %s; want %d", created.Code, created.Body, http.StatusCreated)
	}
	var result project.Project
	if err := json.NewDecoder(created.Body).Decode(&result); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if result.ID == "client-id" || result.ID == "" {
		t.Errorf("created ID = %q, want server-generated ID", result.ID)
	}
	if result.Status != project.Planned {
		t.Errorf("created status = %q, want %q", result.Status, project.Planned)
	}
	if result.Name != "Carva MVP" {
		t.Errorf("created name = %q, want entered name", result.Name)
	}

	lookup := httptest.NewRequest(http.MethodGet, "/api/projects/"+result.ID, nil)
	read := httptest.NewRecorder()
	handler.ServeHTTP(read, lookup)
	if read.Code != http.StatusOK {
		t.Fatalf("GET status = %d, body = %s; want %d", read.Code, read.Body, http.StatusOK)
	}
	var recovered project.Project
	if err := json.NewDecoder(read.Body).Decode(&recovered); err != nil {
		t.Fatalf("decode read response: %v", err)
	}
	if recovered != result {
		t.Errorf("GET = %+v, want %+v", recovered, result)
	}
}

func TestProjectHTTPValidationAndMissingProject(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	handler := httpapi.NewHandler(projects.NewService(carvasqlite.NewProjectRepository(db)), nil, nil)

	invalid := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(`{"name":" \t "}`))
	invalid.Header.Set("Content-Type", "application/json")
	invalidResponse := httptest.NewRecorder()
	handler.ServeHTTP(invalidResponse, invalid)
	if invalidResponse.Code != http.StatusBadRequest {
		t.Errorf("invalid POST status = %d, want %d", invalidResponse.Code, http.StatusBadRequest)
	}

	missing := httptest.NewRequest(http.MethodGet, "/api/projects/missing", nil)
	missingResponse := httptest.NewRecorder()
	handler.ServeHTTP(missingResponse, missing)
	if missingResponse.Code != http.StatusNotFound {
		t.Errorf("missing GET status = %d, want %d", missingResponse.Code, http.StatusNotFound)
	}
}

func TestProjectHTTPDoesNotReportPersistenceFailureAsSuccess(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	if err := db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	handler := httpapi.NewHandler(projects.NewService(carvasqlite.NewProjectRepository(db)), nil, nil)
	request := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(`{"name":"Carva MVP"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Errorf("POST after storage failure status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}
