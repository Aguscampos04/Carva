package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"io"
	"net/http"
	"strconv"

	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	effortdomain "github.com/Aguscampos04/Carva/internal/domain/effort"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

// SprintService contains the application operations required by the Sprint API.
type SprintService interface {
	Create(context.Context, string, string) (sprint.Sprint, error)
	Get(context.Context, string, int64) (sprint.Sprint, error)
	ListClosed(context.Context, string) ([]sprint.Sprint, error)
	Close(context.Context, string, int64) (sprint.Sprint, error)
	AssignBacklogItem(context.Context, string, int64, string) (sprint.PlannedBacklogItem, error)
	ListPlannedBacklogItems(context.Context, string, int64) ([]sprint.PlannedBacklogItem, error)
	CompleteBacklogItem(context.Context, string, int64, string, bool) (sprint.PlannedBacklogItem, error)
	CompleteBacklogItemWithEffort(context.Context, string, int64, string, bool, *effortdomain.Input) (sprint.PlannedBacklogItem, error)
}

type createSprintRequest struct {
	Goal string `json:"goal"`
}

type assignBacklogItemRequest struct {
	BacklogItemID string `json:"backlog_item_id"`
}

type completeBacklogItemRequest struct {
	AcceptanceCriteriaVerified *bool                         `json:"acceptance_criteria_verified"`
	Effort                     *completeBacklogEffortRequest `json:"effort,omitempty"`
}

type completeBacklogEffortRequest struct {
	MemberID string       `json:"member_id"`
	WorkDate project.Date `json:"work_date"`
	Activity string       `json:"activity"`
	Hours    hours.Value  `json:"hours"`
}

func registerSprintRoutes(mux *http.ServeMux, service SprintService) {
	mux.HandleFunc("GET /api/projects/{projectID}/sprints/closed", func(w http.ResponseWriter, r *http.Request) {
		values, err := service.ListClosed(r.Context(), r.PathValue("projectID"))
		if err != nil {
			writeSprintError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, values)
	})

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

	mux.HandleFunc("POST /api/projects/{projectID}/sprints/{sprintID}/close", func(w http.ResponseWriter, r *http.Request) {
		sprintID, err := parseSprintID(r.PathValue("sprintID"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "Sprint ID must be a positive integer")
			return
		}
		closed, err := service.Close(r.Context(), r.PathValue("projectID"), sprintID)
		if err != nil {
			writeSprintError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, closed)
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

	mux.HandleFunc("POST /api/projects/{projectID}/sprints/{sprintID}/backlog-items/{backlogItemID}/completion", func(w http.ResponseWriter, r *http.Request) {
		sprintID, err := parseSprintID(r.PathValue("sprintID"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "Sprint ID must be a positive integer")
			return
		}
		var input completeBacklogItemRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&input); err != nil || input.AcceptanceCriteriaVerified == nil {
			writeError(w, http.StatusBadRequest, "acceptance_criteria_verified is required")
			return
		}
		var trailing any
		if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, "request must contain one JSON value")
			return
		}
		var effortInput *effortdomain.Input
		if input.Effort != nil {
			effortInput = &effortdomain.Input{
				MemberID: input.Effort.MemberID, StoryID: r.PathValue("backlogItemID"),
				WorkDate: input.Effort.WorkDate, Activity: input.Effort.Activity, Hours: input.Effort.Hours,
			}
		}
		completed, err := service.CompleteBacklogItemWithEffort(r.Context(), r.PathValue("projectID"), sprintID, r.PathValue("backlogItemID"), *input.AcceptanceCriteriaVerified, effortInput)
		if err != nil {
			writeSprintError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, completed)
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
	case errors.Is(err, project.ErrFinalized):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, sprint.ErrNextSprintRequired), errors.Is(err, sprint.ErrAlreadyClosed), errors.Is(err, sprint.ErrClosed):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, projects.ErrBusy):
		writeError(w, http.StatusServiceUnavailable, projects.ErrBusy.Error())
	case errors.Is(err, sprint.ErrInvalidProjectID), errors.Is(err, sprint.ErrInvalidGoal), errors.Is(err, sprint.ErrInvalidBacklogItemID):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, project.ErrNotFound), errors.Is(err, sprint.ErrNotFound), errors.Is(err, backlog.ErrNotFound):
		writeError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, sprint.ErrBacklogItemAlreadyAssigned):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, sprint.ErrEffortRequired):
		writeErrorWithCode(w, http.StatusConflict, "effort_required", "Registrá horas trabajadas para completar esta historia")
	case errors.Is(err, sprint.ErrEffortAlreadyRegistered):
		writeErrorWithCode(w, http.StatusConflict, "effort_already_registered", "La historia ya tiene horas activas. Actualizá la vista antes de reintentar")
	case errors.Is(err, backlog.ErrAcceptanceCriteriaUnverified):
		writeError(w, http.StatusConflict, "acceptance criteria must be verified before completing the story")
	case errors.Is(err, effortdomain.ErrDateRange), errors.Is(err, effortdomain.ErrUnknownEnd), errors.Is(err, effortdomain.ErrDailyLimit), errors.Is(err, effortdomain.ErrSprint):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, effortdomain.ErrInvalid), errors.Is(err, hours.ErrInvalid):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, effortdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, "effort member not found in project")
	case errors.Is(err, sprint.ErrBacklogItemNotAssigned):
		writeError(w, http.StatusNotFound, "resource not found")
	default:
		writeError(w, http.StatusInternalServerError, "could not handle Sprint request")
	}
}
