package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	app "github.com/Aguscampos04/Carva/internal/application/defects"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	domain "github.com/Aguscampos04/Carva/internal/domain/defect"
)

type DefectService interface {
	Create(context.Context, string, string, domain.Input) (domain.Defect, error)
	ListByStory(context.Context, string, string) ([]domain.Defect, error)
}

func decodeDefectInput(body io.Reader) (domain.Input, error) {
	fields, err := decodeProjectObject(body, "description", "severity", "status", "detection_sprint_id", "resolution_sprint_id")
	if err != nil {
		return domain.Input{}, err
	}
	for _, key := range []string{"description", "severity", "status"} {
		raw, exists := fields[key]
		if !exists || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return domain.Input{}, domain.ErrInvalidDescription
		}
	}
	data, err := json.Marshal(fields)
	if err != nil {
		return domain.Input{}, err
	}
	var input domain.Input
	if err := json.Unmarshal(data, &input); err != nil {
		return domain.Input{}, err
	}
	return input, input.Validate()
}

func writeDefectError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, domain.ErrInvalidIdentity), errors.Is(err, domain.ErrInvalidDescription),
		errors.Is(err, domain.ErrInvalidSeverity), errors.Is(err, domain.ErrInvalidStatus), errors.Is(err, domain.ErrInvalidSprint):
		writeError(w, http.StatusBadRequest, "invalid defect request")
	case errors.Is(err, domain.ErrProjectNotFound), errors.Is(err, domain.ErrBacklogItemNotFound), errors.Is(err, domain.ErrSprintNotFound):
		writeError(w, http.StatusNotFound, "defect project, story, or Sprint not found")
	case errors.Is(err, projects.ErrBusy):
		writeError(w, http.StatusServiceUnavailable, projects.ErrBusy.Error())
	default:
		writeError(w, http.StatusInternalServerError, "could not process defect request")
	}
	return true
}

// WithDefects adds US-18 create/list endpoints without changing existing API constructors.
func WithDefects(base http.Handler, service *app.Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", base)
	story := "/api/projects/{projectID}/backlog/{itemID}/defects"
	mux.HandleFunc("POST "+story, func(w http.ResponseWriter, r *http.Request) {
		input, err := decodeDefectInput(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid defect request")
			return
		}
		created, err := service.Create(r.Context(), r.PathValue("projectID"), r.PathValue("itemID"), input)
		if !writeDefectError(w, err) {
			writeJSON(w, http.StatusCreated, created)
		}
	})
	mux.HandleFunc("GET "+story, func(w http.ResponseWriter, r *http.Request) {
		items, err := service.ListByStory(r.Context(), r.PathValue("projectID"), r.PathValue("itemID"))
		if !writeDefectError(w, err) {
			writeJSON(w, http.StatusOK, items)
		}
	})
	return mux
}
