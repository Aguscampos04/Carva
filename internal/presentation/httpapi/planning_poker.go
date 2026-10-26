package httpapi

import (
	"context"
	"errors"
	"net/http"
	"strings"

	backlogdomain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	memberdomain "github.com/Aguscampos04/Carva/internal/domain/member"
	domain "github.com/Aguscampos04/Carva/internal/domain/planningpoker"
)

// PlanningPokerService exposes project-scoped round and vote operations.
type PlanningPokerService interface {
	StartRound(context.Context, string, string) (domain.Round, error)
	ListRounds(context.Context, string, string) ([]domain.Round, error)
	RecordVote(context.Context, string, string, string, string, int) error
	VisibleVotes(context.Context, string, string, string, string) ([]domain.Vote, error)
	FinalizeRound(context.Context, string, string, string, string) (domain.RoundResult, error)
	RoundResult(context.Context, string, string, string, string) (domain.RoundResult, error)
	RegisterAgreement(context.Context, string, string, string, string) (int, error)
}

type planningPokerVoteRequest struct {
	StoryPoints *int `json:"story_points"`
}

func registerPlanningPokerRoutes(mux *http.ServeMux, service PlanningPokerService) {
	base := "/api/projects/{projectID}/backlog/{storyID}/estimation"
	mux.HandleFunc("POST "+base+"/rounds", func(w http.ResponseWriter, r *http.Request) {
		round, err := service.StartRound(r.Context(), r.PathValue("projectID"), r.PathValue("storyID"))
		if planningPokerError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, round)
	})
	mux.HandleFunc("GET "+base+"/rounds", func(w http.ResponseWriter, r *http.Request) {
		rounds, err := service.ListRounds(r.Context(), r.PathValue("projectID"), r.PathValue("storyID"))
		if planningPokerError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, rounds)
	})
	mux.HandleFunc("PUT "+base+"/rounds/{roundID}/votes/{memberID}", func(w http.ResponseWriter, r *http.Request) {
		var input planningPokerVoteRequest
		if err := decodeSingleJSON(r, &input); err != nil || input.StoryPoints == nil {
			writeError(w, http.StatusBadRequest, "invalid JSON request")
			return
		}
		err := service.RecordVote(r.Context(), r.PathValue("projectID"), r.PathValue("storyID"), r.PathValue("roundID"), r.PathValue("memberID"), *input.StoryPoints)
		if planningPokerError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"saved": true})
	})
	mux.HandleFunc("GET "+base+"/rounds/{roundID}/votes", func(w http.ResponseWriter, r *http.Request) {
		memberID := r.URL.Query().Get("member_id")
		if strings.TrimSpace(memberID) == "" {
			writeError(w, http.StatusBadRequest, "member_id is required")
			return
		}
		votes, err := service.VisibleVotes(r.Context(), r.PathValue("projectID"), r.PathValue("storyID"), r.PathValue("roundID"), memberID)
		if planningPokerError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, votes)
	})
	mux.HandleFunc("POST "+base+"/rounds/{roundID}/finalize", func(w http.ResponseWriter, r *http.Request) {
		memberID := r.URL.Query().Get("member_id")
		if strings.TrimSpace(memberID) == "" {
			writeError(w, http.StatusBadRequest, "member_id is required")
			return
		}
		result, err := service.FinalizeRound(r.Context(), r.PathValue("projectID"), r.PathValue("storyID"), r.PathValue("roundID"), memberID)
		if planningPokerError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("GET "+base+"/rounds/{roundID}/result", func(w http.ResponseWriter, r *http.Request) {
		memberID := r.URL.Query().Get("member_id")
		if strings.TrimSpace(memberID) == "" {
			writeError(w, http.StatusBadRequest, "member_id is required")
			return
		}
		result, err := service.RoundResult(r.Context(), r.PathValue("projectID"), r.PathValue("storyID"), r.PathValue("roundID"), memberID)
		if planningPokerError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	mux.HandleFunc("POST "+base+"/rounds/{roundID}/agreement", func(w http.ResponseWriter, r *http.Request) {
		memberID := r.URL.Query().Get("member_id")
		if strings.TrimSpace(memberID) == "" {
			writeError(w, http.StatusBadRequest, "member_id is required")
			return
		}
		points, err := service.RegisterAgreement(r.Context(), r.PathValue("projectID"), r.PathValue("storyID"), r.PathValue("roundID"), memberID)
		if planningPokerError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"story_points": points})
	})
}

func planningPokerError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, domain.ErrInvalidMemberID), errors.Is(err, domain.ErrInvalidStoryID),
		errors.Is(err, domain.ErrInvalidRoundID), errors.Is(err, domain.ErrInvalidProjectID),
		errors.Is(err, domain.ErrInvalidStoryPoints), errors.Is(err, domain.ErrInvalidRoundNumber),
		errors.Is(err, domain.ErrDuplicateParticipant), errors.Is(err, domain.ErrNotParticipant):
		writeError(w, http.StatusBadRequest, "invalid Planning Poker participant or vote")
	case errors.Is(err, domain.ErrNotFound), errors.Is(err, domain.ErrRoundNotFound),
		errors.Is(err, backlogdomain.ErrNotFound), errors.Is(err, backlogdomain.ErrProjectNotFound),
		errors.Is(err, memberdomain.ErrProjectNotFound):
		writeError(w, http.StatusNotFound, "Planning Poker project, story, or round not found")
	case errors.Is(err, domain.ErrVotesPending):
		writeError(w, http.StatusConflict, "faltan votos de participantes para finalizar la ronda")
	case errors.Is(err, domain.ErrNoConsensus):
		writeError(w, http.StatusConflict, "la ronda no tiene una estimación unánime")
	case errors.Is(err, domain.ErrRoundNotFinalized):
		writeError(w, http.StatusConflict, "la ronda todavía no está finalizada")
	case errors.Is(err, domain.ErrRoundExists), errors.Is(err, domain.ErrRoundNotOpen):
		writeError(w, http.StatusConflict, "Planning Poker round cannot accept this operation")
	case errors.Is(err, domain.ErrNoParticipants):
		writeError(w, http.StatusUnprocessableEntity, "project must have registered members to start a round")
	default:
		writeError(w, http.StatusInternalServerError, "could not process Planning Poker request")
	}
	return true
}
