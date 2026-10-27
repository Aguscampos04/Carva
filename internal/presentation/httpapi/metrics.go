package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	appmetrics "github.com/Aguscampos04/Carva/internal/application/metrics"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

// WithMetrics adds US-20 read-only metric endpoints without changing older constructors.
func WithMetrics(base http.Handler, service *appmetrics.Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", base)
	mux.HandleFunc("GET /api/projects/{projectID}/metrics/story-points", func(w http.ResponseWriter, r *http.Request) {
		result, err := service.ProjectStoryPoints(r.Context(), r.PathValue("projectID"))
		if !writeMetricsError(w, err) {
			writeJSON(w, http.StatusOK, result)
		}
	})
	mux.HandleFunc("GET /api/projects/{projectID}/sprints/{sprintID}/metrics/story-points", func(w http.ResponseWriter, r *http.Request) {
		sprintID, err := strconv.ParseInt(r.PathValue("sprintID"), 10, 64)
		if err != nil || sprintID < 1 {
			writeError(w, http.StatusBadRequest, "Sprint ID must be a positive integer")
			return
		}
		result, err := service.SprintStoryPoints(r.Context(), r.PathValue("projectID"), sprintID)
		if !writeMetricsError(w, err) {
			writeJSON(w, http.StatusOK, result)
		}
	})
	return mux
}

func writeMetricsError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, project.ErrNotFound) || errors.Is(err, sprint.ErrNotFound) {
		writeError(w, http.StatusNotFound, "metric resource not found")
	} else {
		writeError(w, http.StatusInternalServerError, "could not retrieve project metrics")
	}
	return true
}
