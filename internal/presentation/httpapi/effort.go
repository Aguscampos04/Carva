package httpapi

import (
	"encoding/json"
	"errors"
	app "github.com/Aguscampos04/Carva/internal/application/effort"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	domain "github.com/Aguscampos04/Carva/internal/domain/effort"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"net/http"
)

func decodeEffort(r *http.Request, create bool) (domain.Input, error) {
	keys := []string{"member_id", "work_date", "activity", "hours", "sprint_id"}
	if !create {
		keys = append(keys, "backlog_item_id")
	}
	fields, err := decodeProjectObject(r.Body, keys...)
	if err != nil {
		return domain.Input{}, err
	}
	for _, key := range keys {
		if key != "sprint_id" && (len(fields[key]) == 0 || string(fields[key]) == "null") {
			return domain.Input{}, domain.ErrInvalid
		}
	}
	data, err := json.Marshal(fields)
	if err != nil {
		return domain.Input{}, err
	}
	var in domain.Input
	if err = json.Unmarshal(data, &in); err != nil {
		return in, err
	}
	_, in.HasSprint = fields["sprint_id"]
	if create {
		in.StoryID = r.PathValue("itemID")
	}
	return in, in.Validate()
}
func writeEffortError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	code := 500
	message := "could not process effort request"
	switch {
	case errors.Is(err, domain.ErrInvalid):
		code = 400
		message = err.Error()
	case errors.Is(err, domain.ErrNotFound), errors.Is(err, project.ErrNotFound):
		code = 404
		message = "effort resource not found in project"
	case errors.Is(err, domain.ErrSprint), errors.Is(err, domain.ErrDateRange), errors.Is(err, domain.ErrUnknownEnd), errors.Is(err, domain.ErrDailyLimit):
		code = 409
		message = err.Error()
	case errors.Is(err, projects.ErrBusy):
		code = 503
		message = projects.ErrBusy.Error()
	}
	writeError(w, code, message)
	return true
}

// WithEffort adds US-16 to the existing API without changing older constructors.
func WithEffort(base http.Handler, s *app.Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", base)
	story := "/api/projects/{projectID}/backlog/{itemID}"
	entry := "/api/projects/{projectID}/effort/{entryID}"
	mux.HandleFunc("POST "+story+"/effort", func(w http.ResponseWriter, r *http.Request) {
		in, err := decodeEffort(r, true)
		if err != nil {
			writeError(w, 400, "invalid effort request")
			return
		}
		result, err := s.Save(r.Context(), r.PathValue("projectID"), "", in)
		if !writeEffortError(w, err) {
			writeJSON(w, 201, result)
		}
	})
	mux.HandleFunc("PUT "+entry, func(w http.ResponseWriter, r *http.Request) {
		in, err := decodeEffort(r, false)
		if err != nil {
			writeError(w, 400, "invalid effort request")
			return
		}
		result, err := s.Save(r.Context(), r.PathValue("projectID"), r.PathValue("entryID"), in)
		if !writeEffortError(w, err) {
			writeJSON(w, 200, result)
		}
	})
	mux.HandleFunc("GET "+entry, func(w http.ResponseWriter, r *http.Request) {
		v, err := s.Get(r.Context(), r.PathValue("projectID"), r.PathValue("entryID"))
		if !writeEffortError(w, err) {
			writeJSON(w, 200, v)
		}
	})
	mux.HandleFunc("GET "+story+"/effort", func(w http.ResponseWriter, r *http.Request) {
		v, err := s.List(r.Context(), r.PathValue("projectID"), r.PathValue("itemID"))
		if !writeEffortError(w, err) {
			writeJSON(w, 200, v)
		}
	})
	mux.HandleFunc("GET "+story+"/effort-context", func(w http.ResponseWriter, r *http.Request) {
		v, err := s.Context(r.Context(), r.PathValue("projectID"), r.PathValue("itemID"))
		if !writeEffortError(w, err) {
			writeJSON(w, 200, v)
		}
	})
	mux.HandleFunc("POST "+entry+"/delete", func(w http.ResponseWriter, r *http.Request) {
		fields, err := decodeProjectObject(r.Body, "confirmed")
		var confirmed bool
		if err == nil {
			err = json.Unmarshal(fields["confirmed"], &confirmed)
		}
		if err != nil || !confirmed {
			writeError(w, 400, "effort deletion must be confirmed")
			return
		}
		if !writeEffortError(w, s.Delete(r.Context(), r.PathValue("projectID"), r.PathValue("entryID"))) {
			w.WriteHeader(204)
		}
	})
	return mux
}
