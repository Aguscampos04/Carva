package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

// SprintService contains the application operations required by the Sprint API.
type SprintService interface {
	Create(context.Context, string, string) (sprint.Sprint, error)
	Get(context.Context, string, int64) (sprint.Sprint, error)
	AssignBacklogItem(context.Context, string, int64, string) (sprint.PlannedBacklogItem, error)
	ListPlannedBacklogItems(context.Context, string, int64) ([]sprint.PlannedBacklogItem, error)
}

type createSprintRequest struct {
	Goal string `json:"goal"`
}

type assignBacklogItemRequest struct {
	BacklogItemID string `json:"backlog_item_id"`
}

func registerSprintRoutes(mux *http.ServeMux, service SprintService) {
	mux.HandleFunc("POST /api/projects/{projectID}/sprints", func(w http.ResponseWriter, r *http.Request) {
		var input createSprintRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON request")
			return
		}
		var trailing any
		if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, "request must contain one JSON value")
			return
		}

		projectID := r.PathValue("projectID")
		created, err := service.Create(r.Context(), projectID, input.Goal)
		if err != nil {
			writeSprintError(w, err)
			return
		}
		w.Header().Set("Location", "/api/projects/"+projectID+"/sprints/"+strconv.FormatInt(created.ID, 10))
		writeJSON(w, http.StatusCreated, created)
	})

	mux.HandleFunc("GET /api/projects/{projectID}/sprints/{sprintID}", func(w http.ResponseWriter, r *http.Request) {
		sprintID, err := strconv.ParseInt(r.PathValue("sprintID"), 10, 64)
		if err != nil || sprintID < 1 {
			writeError(w, http.StatusBadRequest, "Sprint ID must be a positive integer")
			return
		}
		found, err := service.Get(r.Context(), r.PathValue("projectID"), sprintID)
		if err != nil {
			writeSprintError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, found)
	})

	mux.HandleFunc("POST /api/projects/{projectID}/sprints/{sprintID}/backlog-items", func(w http.ResponseWriter, r *http.Request) {
		sprintID, err := parseSprintID(r.PathValue("sprintID"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "Sprint ID must be a positive integer")
			return
		}
		var input assignBacklogItemRequest
		if err := decodeSingleJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON request")
			return
		}
		created, err := service.AssignBacklogItem(r.Context(), r.PathValue("projectID"), sprintID, input.BacklogItemID)
		if err != nil {
			writeSprintError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	})

	mux.HandleFunc("GET /api/projects/{projectID}/sprints/{sprintID}/backlog-items", func(w http.ResponseWriter, r *http.Request) {
		sprintID, err := parseSprintID(r.PathValue("sprintID"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "Sprint ID must be a positive integer")
			return
		}
		items, err := service.ListPlannedBacklogItems(r.Context(), r.PathValue("projectID"), sprintID)
		if err != nil {
			writeSprintError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
}

func parseSprintID(value string) (int64, error) {
	sprintID, err := strconv.ParseInt(value, 10, 64)
	if err != nil || sprintID < 1 {
		return 0, sprint.ErrNotFound
	}
	return sprintID, nil
}

func writeSprintError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sprint.ErrInvalidProjectID), errors.Is(err, sprint.ErrInvalidGoal), errors.Is(err, sprint.ErrInvalidBacklogItemID):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, project.ErrNotFound), errors.Is(err, sprint.ErrNotFound), errors.Is(err, backlog.ErrNotFound):
		writeError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, sprint.ErrBacklogItemAlreadyAssigned):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "could not handle Sprint request")
	}
}
