package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
)

// BacklogService contains the application operations exposed by the HTTP API.
type BacklogService interface {
	Create(context.Context, string, domain.Input) (domain.Item, error)
	Update(context.Context, string, string, domain.Input) (domain.Item, error)
	List(context.Context, string) ([]domain.Item, error)
}

type backlogErrorResponse struct {
	Error string `json:"error"`
}

func registerBacklogRoutes(mux *http.ServeMux, service BacklogService) {
	mux.HandleFunc("POST /api/projects/{projectID}/backlog", createBacklogItem(service))
	mux.HandleFunc("PUT /api/projects/{projectID}/backlog/{itemID}", updateBacklogItem(service))
	mux.HandleFunc("GET /api/projects/{projectID}/backlog", listBacklogItems(service))
}

func createBacklogItem(service BacklogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input domain.Input
		if err := decodeSingleJSON(r, &input); err != nil {
			writeBacklogError(w, http.StatusBadRequest, "invalid JSON request")
			return
		}
		created, err := service.Create(r.Context(), r.PathValue("projectID"), input)
		if writeBacklogServiceError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

func updateBacklogItem(service BacklogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input domain.Input
		if err := decodeSingleJSON(r, &input); err != nil {
			writeBacklogError(w, http.StatusBadRequest, "invalid JSON request")
			return
		}
		updated, err := service.Update(r.Context(), r.PathValue("projectID"), r.PathValue("itemID"), input)
		if writeBacklogServiceError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, updated)
	}
}

func listBacklogItems(service BacklogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := service.List(r.Context(), r.PathValue("projectID"))
		if writeBacklogServiceError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

func decodeSingleJSON(r *http.Request, destination any) error {
	if err := json.NewDecoder(r.Body).Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := json.NewDecoder(r.Body).Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("request must contain one JSON value")
	}
	return nil
}

func writeBacklogServiceError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, domain.ErrProjectNotFound), errors.Is(err, domain.ErrNotFound):
		writeBacklogError(w, http.StatusNotFound, "backlog project or item not found")
	case errors.Is(err, domain.ErrInvalidID), errors.Is(err, domain.ErrInvalidProjectID),
		errors.Is(err, domain.ErrInvalidTitle), errors.Is(err, domain.ErrInvalidDescription),
		errors.Is(err, domain.ErrInvalidPriority), errors.Is(err, domain.ErrInvalidProvisionalStoryPoints),
		errors.Is(err, domain.ErrInvalidAcceptanceCriteria), errors.Is(err, domain.ErrEstimateFinalized):
		writeBacklogError(w, http.StatusBadRequest, "invalid backlog item")
	default:
		writeBacklogError(w, http.StatusInternalServerError, "could not process Product Backlog request")
	}
	return true
}

func writeBacklogError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, backlogErrorResponse{Error: message})
}
