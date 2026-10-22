package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

// US-04 BDD-01/02/03: status stays queryable across the project lifecycle.
func TestUS04ProjectStatusCanBeConsultedThroughoutLifecycle(t *testing.T) {
	handler := newTestHandler(t)
	projectID := createProject(t, handler, "Carva Demo")
	projectURL := "/api/projects/" + projectID

	assertProjectStatus := func(want string) {
		t.Helper()
		response := performJSON(t, handler, http.MethodGet, projectURL, "")
		if response.Code != http.StatusOK {
			t.Fatalf("GET project status = %d, body = %s; want %d", response.Code, response.Body, http.StatusOK)
		}
		var result struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatalf("decode project response: %v", err)
		}
		if result.Status != want {
			t.Errorf("consulted project state = %q, want %q", result.Status, want)
		}
	}

	assertProjectStatus("Planificado")

	sprint := performJSON(t, handler, http.MethodPost, projectURL+"/sprints", `{"goal":"Entregar el MVP"}`)
	if sprint.Code != http.StatusCreated {
		t.Fatalf("create first Sprint status = %d, body = %s", sprint.Code, sprint.Body)
	}
	assertProjectStatus("En progreso")

	finalized := performJSON(t, handler, http.MethodPost, projectURL+"/finalize", `{"confirmed":true}`)
	if finalized.Code != http.StatusOK {
		t.Fatalf("finalize project status = %d, body = %s", finalized.Code, finalized.Body)
	}
	assertProjectStatus("Finalizado")
}
