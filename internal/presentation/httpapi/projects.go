package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Aguscampos04/Carva/internal/domain/project"
)

// ProjectService contains the application operations required by the HTTP API.
type ProjectService interface {
	Create(ctx context.Context, name string) (project.Project, error)
	Get(ctx context.Context, id string) (project.Project, error)
}

// NewHandler returns the REST endpoints for project creation and consultation.
func NewHandler(service ProjectService, sprintServices ...SprintService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/projects", createProject(service))
	mux.HandleFunc("GET /api/projects/{id}", getProject(service))
	if len(sprintServices) > 0 && sprintServices[0] != nil {
		registerSprintRoutes(mux, sprintServices[0])
	}
	return mux
}

type createProjectRequest struct {
	Name string `json:"name"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func createProject(service ProjectService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input createProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON request")
			return
		}

		var trailing any
		if err := json.NewDecoder(r.Body).Decode(&trailing); !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, "request must contain one JSON value")
			return
		}

		created, err := service.Create(r.Context(), input.Name)
		if errors.Is(err, project.ErrInvalidName) {
			writeError(w, http.StatusBadRequest, "project name must not be blank")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not create project")
			return
		}

		w.Header().Set("Location", "/api/projects/"+created.ID)
		writeJSON(w, http.StatusCreated, created)
	}
}

func getProject(service ProjectService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		found, err := service.Get(r.Context(), r.PathValue("id"))
		if errors.Is(err, project.ErrNotFound) {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not retrieve project")
			return
		}
		writeJSON(w, http.StatusOK, found)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
