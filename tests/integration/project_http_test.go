package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
	carvasqlite "github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

func TestProjectHTTPCreateReadAndRejectClientControlledIdentityOrStatus(t *testing.T) {
	ctx := context.Background()
	db, err := carvasqlite.OpenPath(ctx, filepath.Join(t.TempDir(), "project.db"))
	if err != nil {
		t.Fatalf("OpenPath() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := carvasqlite.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	handler := httpapi.NewHandler(projects.NewService(carvasqlite.NewProjectRepository(db)))

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
	ctx := context.Background()
	db, err := carvasqlite.OpenPath(ctx, filepath.Join(t.TempDir(), "project.db"))
	if err != nil {
		t.Fatalf("OpenPath() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := carvasqlite.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	handler := httpapi.NewHandler(projects.NewService(carvasqlite.NewProjectRepository(db)))

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
	ctx := context.Background()
	db, err := carvasqlite.OpenPath(ctx, filepath.Join(t.TempDir(), "project.db"))
	if err != nil {
		t.Fatalf("OpenPath() error = %v", err)
	}
	if err := carvasqlite.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	handler := httpapi.NewHandler(projects.NewService(carvasqlite.NewProjectRepository(db)))
	request := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(`{"name":"Carva MVP"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Errorf("POST after storage failure status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}
