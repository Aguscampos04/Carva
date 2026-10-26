package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	appplanningpoker "github.com/Aguscampos04/Carva/internal/application/planningpoker"
	domain "github.com/Aguscampos04/Carva/internal/domain/planningpoker"
)

// PlanningPokerRepository persists estimation rounds and hidden votes.
type PlanningPokerRepository struct{ db *sql.DB }

func NewPlanningPokerRepository(db *sql.DB) *PlanningPokerRepository {
	return &PlanningPokerRepository{db: db}
}

func (r *PlanningPokerRepository) CreateRound(ctx context.Context, round domain.Round) error {
	return withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		var storyExists bool
		if err := conn.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM backlog_items WHERE id = ? AND project_id = ?)`, round.StoryID, round.ProjectID,
		).Scan(&storyExists); err != nil {
			return fmt.Errorf("check Planning Poker story: %w", err)
		}
		if !storyExists {
			return domain.ErrNotFound
		}
		var previousNumber int
		var previousStatus string
		err := conn.QueryRowContext(ctx, `
			SELECT round_number, status FROM estimation_rounds
			WHERE story_id = ? ORDER BY round_number DESC LIMIT 1`, round.StoryID,
		).Scan(&previousNumber, &previousStatus)
		expectedNumber := 1
		switch {
		case errors.Is(err, sql.ErrNoRows):
		case err != nil:
			return fmt.Errorf("find previous Planning Poker round: %w", err)
		default:
			if domain.RoundStatus(previousStatus) != domain.RoundFinalized {
				return domain.ErrRoundNotFinalized
			}
			expectedNumber = previousNumber + 1
		}
		if round.Number != expectedNumber {
			return domain.ErrRoundExists
		}
		if _, err := conn.ExecContext(ctx, `
			INSERT INTO estimation_rounds (id, project_id, story_id, round_number, status)
			VALUES (?, ?, ?, ?, ?)`, round.ID, round.ProjectID, round.StoryID, round.Number, round.Status,
		); err != nil {
			return fmt.Errorf("insert Planning Poker round: %w", err)
		}
		for _, memberID := range round.Participants {
			var belongs bool
			if err := conn.QueryRowContext(ctx,
				`SELECT EXISTS(SELECT 1 FROM members WHERE id = ? AND project_id = ?)`, memberID, round.ProjectID,
			).Scan(&belongs); err != nil {
				return fmt.Errorf("check Planning Poker participant: %w", err)
			}
			if !belongs {
				return domain.ErrNotParticipant
			}
			if _, err := conn.ExecContext(ctx,
				`INSERT INTO estimation_round_participants (round_id, member_id) VALUES (?, ?)`, round.ID, memberID,
			); err != nil {
				return fmt.Errorf("insert Planning Poker participant: %w", err)
			}
		}
		return nil
	})
}

func (r *PlanningPokerRepository) ListRounds(ctx context.Context, projectID, storyID string) ([]domain.Round, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, project_id, story_id, round_number, status
		FROM estimation_rounds WHERE project_id = ? AND story_id = ? ORDER BY round_number`, projectID, storyID)
	if err != nil {
		return nil, fmt.Errorf("query Planning Poker rounds: %w", err)
	}
	defer rows.Close()
	var rounds []domain.Round
	for rows.Next() {
		round, err := scanPlanningPokerRound(rows)
		if err != nil {
			return nil, fmt.Errorf("scan Planning Poker round: %w", err)
		}
		participants, err := r.participants(ctx, round.ID)
		if err != nil {
			return nil, err
		}
		round.Participants = participants
		rounds = append(rounds, round)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Planning Poker rounds: %w", err)
	}
	return rounds, nil
}

func (r *PlanningPokerRepository) FindRound(ctx context.Context, projectID, storyID, roundID string) (domain.Round, error) {
	round, err := scanPlanningPokerRound(r.db.QueryRowContext(ctx, `
		SELECT id, project_id, story_id, round_number, status
		FROM estimation_rounds WHERE id = ? AND project_id = ? AND story_id = ?`, roundID, projectID, storyID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Round{}, domain.ErrRoundNotFound
	}
	if err != nil {
		return domain.Round{}, fmt.Errorf("select Planning Poker round: %w", err)
	}
	round.Participants, err = r.participants(ctx, round.ID)
	if err != nil {
		return domain.Round{}, err
	}
	return round, nil
}

func (r *PlanningPokerRepository) UpsertVote(ctx context.Context, round domain.Round, vote domain.Vote) error {
	return withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		stored, err := scanPlanningPokerRound(conn.QueryRowContext(ctx, `
			SELECT id, project_id, story_id, round_number, status
			FROM estimation_rounds WHERE id = ? AND project_id = ? AND story_id = ?`,
			round.ID, round.ProjectID, round.StoryID))
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrRoundNotFound
		}
		if err != nil {
			return fmt.Errorf("select Planning Poker round for vote: %w", err)
		}
		stored.Participants, err = participantsOn(ctx, conn, stored.ID)
		if err != nil {
			return err
		}
		if err := stored.ValidateVote(vote); err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, `
			INSERT INTO estimation_votes (round_id, member_id, story_points)
			VALUES (?, ?, ?)
			ON CONFLICT (round_id, member_id) DO UPDATE SET story_points = excluded.story_points`,
			vote.RoundID, vote.MemberID, vote.StoryPoints)
		if err != nil {
			return fmt.Errorf("upsert Planning Poker vote: %w", err)
		}
		return nil
	})
}

func (r *PlanningPokerRepository) ListVotes(ctx context.Context, roundID string) ([]domain.Vote, error) {
	return votesOn(ctx, r.db, roundID)
}

func (r *PlanningPokerRepository) FinalizeRound(ctx context.Context, projectID, storyID, roundID string) (domain.RoundResult, error) {
	var result domain.RoundResult
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		round, err := scanPlanningPokerRound(conn.QueryRowContext(ctx, `
			SELECT id, project_id, story_id, round_number, status
			FROM estimation_rounds WHERE id = ? AND project_id = ? AND story_id = ?`, roundID, projectID, storyID))
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrRoundNotFound
		}
		if err != nil {
			return fmt.Errorf("select Planning Poker round for finalization: %w", err)
		}
		round.Participants, err = participantsOn(ctx, conn, round.ID)
		if err != nil {
			return err
		}
		votes, err := votesOn(ctx, conn, round.ID)
		if err != nil {
			return err
		}
		result, err = round.Finalize(votes)
		if err != nil {
			return err
		}
		updated, err := conn.ExecContext(ctx,
			`UPDATE estimation_rounds SET status = ? WHERE id = ? AND status = ?`, domain.RoundFinalized, round.ID, domain.RoundOpen)
		if err != nil {
			return fmt.Errorf("finalize Planning Poker round: %w", err)
		}
		rows, err := updated.RowsAffected()
		if err != nil {
			return fmt.Errorf("count finalized Planning Poker rounds: %w", err)
		}
		if rows != 1 {
			return domain.ErrRoundNotOpen
		}
		return nil
	})
	if err != nil {
		return domain.RoundResult{}, err
	}
	return result, nil
}

func (r *PlanningPokerRepository) GetRoundResult(ctx context.Context, projectID, storyID, roundID string) (domain.RoundResult, error) {
	round, err := r.FindRound(ctx, projectID, storyID, roundID)
	if err != nil {
		return domain.RoundResult{}, err
	}
	votes, err := r.ListVotes(ctx, roundID)
	if err != nil {
		return domain.RoundResult{}, err
	}
	return round.Result(votes)
}

func votesOn(ctx context.Context, queryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, roundID string) ([]domain.Vote, error) {
	rows, err := queryer.QueryContext(ctx, `
		SELECT votes.member_id, rounds.story_id, votes.story_points
		FROM estimation_votes AS votes
		JOIN estimation_rounds AS rounds ON rounds.id = votes.round_id
		WHERE votes.round_id = ? ORDER BY votes.member_id`, roundID)
	if err != nil {
		return nil, fmt.Errorf("query Planning Poker votes: %w", err)
	}
	defer rows.Close()
	votes := make([]domain.Vote, 0)
	for rows.Next() {
		var vote domain.Vote
		vote.RoundID = roundID
		if err := rows.Scan(&vote.MemberID, &vote.StoryID, &vote.StoryPoints); err != nil {
			return nil, fmt.Errorf("scan Planning Poker vote: %w", err)
		}
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Planning Poker votes: %w", err)
	}
	return votes, nil
}

type planningPokerScanner interface{ Scan(...any) error }

func scanPlanningPokerRound(scanner planningPokerScanner) (domain.Round, error) {
	var round domain.Round
	var status string
	if err := scanner.Scan(&round.ID, &round.ProjectID, &round.StoryID, &round.Number, &status); err != nil {
		return domain.Round{}, err
	}
	round.Status = domain.RoundStatus(status)
	return round, nil
}

func (r *PlanningPokerRepository) participants(ctx context.Context, roundID string) ([]string, error) {
	return participantsOn(ctx, r.db, roundID)
}

func participantsOn(ctx context.Context, queryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, roundID string) ([]string, error) {
	rows, err := queryer.QueryContext(ctx,
		`SELECT member_id FROM estimation_round_participants WHERE round_id = ? ORDER BY member_id`, roundID)
	if err != nil {
		return nil, fmt.Errorf("query Planning Poker participants: %w", err)
	}
	defer rows.Close()
	participants := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan Planning Poker participant: %w", err)
		}
		participants = append(participants, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Planning Poker participants: %w", err)
	}
	return participants, nil
}

var _ appplanningpoker.Repository = (*PlanningPokerRepository)(nil)
