package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type ProjectService interface {
	Create(context.Context, string) (project.Project, error)
	Get(context.Context, int64) (project.Project, error)
}

type SprintService interface {
	Create(context.Context, int64, string) (sprint.Sprint, error)
	Get(context.Context, int64, int64) (sprint.Sprint, error)
}

type Handler struct {
	projects ProjectService
	sprints  SprintService
}

func NewHandler(projects ProjectService, sprints SprintService) http.Handler {
	h := Handler{projects: projects, sprints: sprints}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/projects", h.createProject)
	mux.HandleFunc("GET /api/projects/{projectID}", h.getProject)
	mux.HandleFunc("POST /api/projects/{projectID}/sprints", h.createSprint)
	mux.HandleFunc("GET /api/projects/{projectID}/sprints/{sprintID}", h.getSprint)
	return mux
}

type projectInput struct {
	Name string `json:"name"`
}

type sprintInput struct {
	Goal string `json:"goal"`
}

func (h Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var input projectInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}
	value, err := h.projects.Create(r.Context(), input.Name)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, value)
}

func (h Handler) getProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "project ID must be a positive integer")
		return
	}
	value, err := h.projects.Get(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func (h Handler) createSprint(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseID(r.PathValue("projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "project ID must be a positive integer")
		return
	}
	var input sprintInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}
	value, err := h.sprints.Create(r.Context(), projectID, input.Goal)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, value)
}

func (h Handler) getSprint(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseID(r.PathValue("projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "project ID must be a positive integer")
		return
	}
	sprintID, err := parseID(r.PathValue("sprintID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Sprint ID must be a positive integer")
		return
	}
	value, err := h.sprints.Get(r.Context(), projectID, sprintID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid ID %q", raw)
	}
	return id, nil
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return errors.New("request must contain exactly one JSON value")
	}
	return nil
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, project.ErrInvalidName), errors.Is(err, project.ErrInvalidID),
		errors.Is(err, sprint.ErrInvalidProjectID), errors.Is(err, sprint.ErrInvalidGoal):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, project.ErrNotFound), errors.Is(err, sprint.ErrNotFound):
		writeError(w, http.StatusNotFound, "resource not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
