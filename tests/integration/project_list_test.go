package integration_test

import (
	"context"
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

func TestProjectListHTTPEmptyThenExistingProjects(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	service := projects.NewService(carvasqlite.NewProjectRepository(db))
	handler := httpapi.NewHandler(service, nil, nil)
	read := func() *httptest.ResponseRecorder {
		t.Helper()
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/projects", nil))
		return response
	}
	empty := read()
	if empty.Code != 200 || strings.TrimSpace(empty.Body.String()) != "[]" {
		t.Fatalf("empty list: %d %s", empty.Code, empty.Body)
	}
	for _, name := range []string{"Zulu", "Alpha", "Alpha"} {
		if _, err := service.Create(context.Background(), name); err != nil {
			t.Fatal(err)
		}
	}
	response := read()
	if response.Code != 200 {
		t.Fatalf("list status %d: %s", response.Code, response.Body)
	}
	var result []project.Project
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 || result[0].Name != "Alpha" || result[1].Name != "Alpha" || result[2].Name != "Zulu" {
		t.Fatalf("unexpected order: %+v", result)
	}
	if result[0].ID >= result[1].ID {
		t.Fatal("homonyms must retain distinct identities in stable order")
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	failed := read()
	if failed.Code != 500 || strings.Contains(failed.Body.String(), "sql:") {
		t.Fatalf("unsafe storage failure: %d %s", failed.Code, failed.Body)
	}
}
