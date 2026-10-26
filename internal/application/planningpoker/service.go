// Package planningpoker coordinates estimation-round and vote use cases.
package planningpoker

import (
	"context"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/member"
	domain "github.com/Aguscampos04/Carva/internal/domain/planningpoker"
	"github.com/google/uuid"
)

// Repository persists rounds and their votes.
type Repository interface {
	CreateInitialRound(context.Context, domain.Round) error
	ListRounds(context.Context, string, string) ([]domain.Round, error)
	FindRound(context.Context, string, string, string) (domain.Round, error)
	UpsertVote(context.Context, domain.Round, domain.Vote) error
	ListVotes(context.Context, string) ([]domain.Vote, error)
	FinalizeRound(context.Context, string, string, string) (domain.RoundResult, error)
	GetRoundResult(context.Context, string, string, string) (domain.RoundResult, error)
}

// BacklogReader verifies that the story belongs to the requested project.
type BacklogReader interface {
	FindByID(context.Context, string, string) (backlog.Item, error)
}

// MemberReader supplies the project's participants when the round opens.
type MemberReader interface {
	List(context.Context, string) ([]member.Member, error)
}

// Service coordinates Planning Poker operations.
type Service struct {
	repository Repository
	backlog    BacklogReader
	members    MemberReader
}

func NewService(repository Repository, backlog BacklogReader, members MemberReader) Service {
	return Service{repository: repository, backlog: backlog, members: members}
}

// StartInitialRound snapshots the registered members for the first round.
func (s Service) StartInitialRound(ctx context.Context, projectID, storyID string) (domain.Round, error) {
	if _, err := s.backlog.FindByID(ctx, projectID, storyID); err != nil {
		return domain.Round{}, fmt.Errorf("find Planning Poker story: %w", err)
	}
	projectMembers, err := s.members.List(ctx, projectID)
	if err != nil {
		return domain.Round{}, fmt.Errorf("list Planning Poker participants: %w", err)
	}
	participantIDs := make([]string, 0, len(projectMembers))
	for _, participant := range projectMembers {
		participantIDs = append(participantIDs, participant.ID)
	}
	round, err := domain.NewRound(uuid.NewString(), projectID, storyID, 1, participantIDs)
	if err != nil {
		return domain.Round{}, err
	}
	if err := s.repository.CreateInitialRound(ctx, round); err != nil {
		return domain.Round{}, fmt.Errorf("create initial Planning Poker round: %w", err)
	}
	return round, nil
}

func (s Service) ListRounds(ctx context.Context, projectID, storyID string) ([]domain.Round, error) {
	if _, err := s.backlog.FindByID(ctx, projectID, storyID); err != nil {
		return nil, fmt.Errorf("find Planning Poker story: %w", err)
	}
	rounds, err := s.repository.ListRounds(ctx, projectID, storyID)
	if err != nil {
		return nil, fmt.Errorf("list Planning Poker rounds: %w", err)
	}
	if rounds == nil {
		return []domain.Round{}, nil
	}
	return rounds, nil
}

func (s Service) RecordVote(ctx context.Context, projectID, storyID, roundID, memberID string, points int) error {
	vote, err := domain.NewVote(memberID, storyID, roundID, points)
	if err != nil {
		return err
	}
	round, err := s.repository.FindRound(ctx, projectID, storyID, roundID)
	if err != nil {
		return fmt.Errorf("find Planning Poker round: %w", err)
	}
	if err := round.ValidateVote(vote); err != nil {
		return err
	}
	if err := s.repository.UpsertVote(ctx, round, vote); err != nil {
		return fmt.Errorf("record Planning Poker vote: %w", err)
	}
	return nil
}

func (s Service) VisibleVotes(ctx context.Context, projectID, storyID, roundID, memberID string) ([]domain.Vote, error) {
	round, err := s.repository.FindRound(ctx, projectID, storyID, roundID)
	if err != nil {
		return nil, fmt.Errorf("find Planning Poker round: %w", err)
	}
	votes, err := s.repository.ListVotes(ctx, roundID)
	if err != nil {
		return nil, fmt.Errorf("list Planning Poker votes: %w", err)
	}
	return round.VisibleVotes(memberID, votes)
}

// FinalizeRound reveals the result to a participant after every participant voted.
func (s Service) FinalizeRound(ctx context.Context, projectID, storyID, roundID, memberID string) (domain.RoundResult, error) {
	round, err := s.repository.FindRound(ctx, projectID, storyID, roundID)
	if err != nil {
		return domain.RoundResult{}, fmt.Errorf("find Planning Poker round for finalization: %w", err)
	}
	if _, err := round.VisibleVotes(memberID, nil); err != nil {
		return domain.RoundResult{}, err
	}
	result, err := s.repository.FinalizeRound(ctx, projectID, storyID, roundID)
	if err != nil {
		return domain.RoundResult{}, fmt.Errorf("finalize Planning Poker round: %w", err)
	}
	return result, nil
}

// RoundResult returns revealed votes to a participant after finalization.
func (s Service) RoundResult(ctx context.Context, projectID, storyID, roundID, memberID string) (domain.RoundResult, error) {
	round, err := s.repository.FindRound(ctx, projectID, storyID, roundID)
	if err != nil {
		return domain.RoundResult{}, fmt.Errorf("find Planning Poker round result: %w", err)
	}
	if _, err := round.VisibleVotes(memberID, nil); err != nil {
		return domain.RoundResult{}, err
	}
	result, err := s.repository.GetRoundResult(ctx, projectID, storyID, roundID)
	if err != nil {
		return domain.RoundResult{}, fmt.Errorf("get Planning Poker round result: %w", err)
	}
	return result, nil
}
