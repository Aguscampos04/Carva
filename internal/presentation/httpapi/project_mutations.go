package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/domain/project"
)

// decodeProjectObject validates the full object with exact, unique, allowed keys.
// It is deliberately limited to new US-02 writes, preserving other contracts.
func decodeProjectObject(body io.Reader, allowed ...string) (map[string]json.RawMessage, error) {
	invalid := errors.New("invalid JSON request")
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(data) {
		return nil, invalid
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil, invalid
	}
	fields := map[string]json.RawMessage{}
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return nil, invalid
		}
		key, ok := token.(string)
		if !ok {
			return nil, invalid
		}
		if _, exists := fields[key]; exists {
			return nil, invalid
		}
		accepted := false
		for _, candidate := range allowed {
			if candidate == key {
				accepted = true
				break
			}
		}
		if !accepted {
			return nil, invalid
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return nil, invalid
		}
		fields[key] = value
	}
	if token, err = decoder.Token(); err != nil || token != json.Delim('}') {
		return nil, invalid
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, invalid
	}
	return fields, nil
}

func decodeProjectPatch(body io.Reader) (project.Patch, error) {
	fields, err := decodeProjectObject(body, "name", "start_date")
	if err != nil {
		return project.Patch{}, err
	}
	var change project.Patch
	if raw, ok := fields["name"]; ok {
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return change, project.ErrInvalidName
		}
		var name string
		if err := json.Unmarshal(raw, &name); err != nil {
			return change, err
		}
		change.Name = &name
	}
	if raw, ok := fields["start_date"]; ok {
		change.HasStartDate = true
		if err := json.Unmarshal(raw, &change.StartDate); err != nil {
			return change, err
		}
	}
	return change, change.Validate()
}

func registerProjectMutationRoutes(mux *http.ServeMux, service ProjectService) {
	mux.HandleFunc("PATCH /api/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		change, err := decodeProjectPatch(r.Body)
		if err != nil {
			writeError(w, 400, "invalid JSON request")
			return
		}
		result, err := service.Update(r.Context(), r.PathValue("id"), change)
		if writeProjectMutationError(w, err, "could not update project") {
			return
		}
		writeJSON(w, 200, result)
	})
	mux.HandleFunc("POST /api/projects/{id}/finalize", func(w http.ResponseWriter, r *http.Request) {
		fields, err := decodeProjectObject(r.Body, "confirmed")
		confirmed := false
		if err == nil {
			err = json.Unmarshal(fields["confirmed"], &confirmed)
		}
		if err != nil || !confirmed {
			writeError(w, 400, "project finalization must be confirmed")
			return
		}
		result, err := service.Finalize(r.Context(), r.PathValue("id"))
		if writeProjectMutationError(w, err, "could not finalize project") {
			return
		}
		writeJSON(w, 200, result)
	})
}

func writeProjectMutationError(w http.ResponseWriter, err error, fallback string) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, project.ErrInvalidName), errors.Is(err, project.ErrInvalidDate), errors.Is(err, project.ErrDateRange), errors.Is(err, project.ErrEmptyPatch):
		writeError(w, 400, err.Error())
	case errors.Is(err, project.ErrNotFound):
		writeError(w, 404, "project not found")
	case errors.Is(err, project.ErrStartAfterToday), errors.Is(err, project.ErrFinalized):
		writeError(w, 409, err.Error())
	case errors.Is(err, projects.ErrBusy):
		writeError(w, 503, projects.ErrBusy.Error())
	default:
		writeError(w, 500, fallback)
	}
	return true
}
