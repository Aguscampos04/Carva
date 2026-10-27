package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type SprintRepository struct{ db *sql.DB }

func NewSprintRepository(db *sql.DB) *SprintRepository {
	return &SprintRepository{db: db}
}

func (r *SprintRepository) GetProject(ctx context.Context, id string) (project.Project, error) {
	return findProject(ctx, r.db, id)
}

func (r *SprintRepository) Create(ctx context.Context, value sprint.Sprint, transition func(project.Project) (project.Project, error)) (sprint.Sprint, error) {
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		current, err := findProject(ctx, conn, value.ProjectID)
		if err != nil {
			return err
		}
		next, err := transition(current)
		if err != nil {
			return err
		}
		if next.Status != current.Status {
			if _, err = conn.ExecContext(ctx, `UPDATE projects SET status=? WHERE id=?`, next.Status, value.ProjectID); err != nil {
				return err
			}
		}
		result, err := conn.ExecContext(ctx, `INSERT INTO sprints (project_id,goal) VALUES (?,?)`, value.ProjectID, value.Goal)
		if err != nil {
			return err
		}
		value.ID, err = result.LastInsertId()
		return err
	})
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("create Sprint: %w", err)
	}
	return value, nil
}

func (r *SprintRepository) Get(ctx context.Context, projectID string, sprintID int64) (sprint.Sprint, error) {
	var value sprint.Sprint
	err := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, goal, closed FROM sprints WHERE project_id = ? AND id = ?`, projectID, sprintID,
	).Scan(&value.ID, &value.ProjectID, &value.Goal, &value.Closed)
	if errors.Is(err, sql.ErrNoRows) {
		return sprint.Sprint{}, sprint.ErrNotFound
	}
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("get Sprint: %w", err)
	}
	return value, nil
}

// Close atomically transfers every incomplete planning snapshot and closes its Sprint.
func (r *SprintRepository) Close(ctx context.Context, projectID string, sprintID int64) (sprint.Sprint, error) {
	var closed sprint.Sprint
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		if err := conn.QueryRowContext(ctx, `
			SELECT id, project_id, goal, closed
			FROM sprints WHERE project_id = ? AND id = ?`, projectID, sprintID,
		).Scan(&closed.ID, &closed.ProjectID, &closed.Goal, &closed.Closed); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return sprint.ErrNotFound
			}
			return fmt.Errorf("get Sprint to close: %w", err)
		}
		closedState, err := closed.Close()
		if err != nil {
			return err
		}
		closed = closedState

		var incompleteCount int
		if err := conn.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM sprint_backlog_items
			WHERE sprint_id = ? AND completion_verified = 0`, sprintID,
		).Scan(&incompleteCount); err != nil {
			return fmt.Errorf("count incomplete Sprint stories: %w", err)
		}

		if incompleteCount > 0 {
			var nextSprintID int64
			err := conn.QueryRowContext(ctx, `
				SELECT id FROM sprints
				WHERE project_id = ? AND id > ? AND closed = 0
				ORDER BY id LIMIT 1`, projectID, sprintID,
			).Scan(&nextSprintID)
			if errors.Is(err, sql.ErrNoRows) {
				return sprint.ErrNextSprintRequired
			}
			if err != nil {
				return fmt.Errorf("find next open Sprint: %w", err)
			}

			_, err = conn.ExecContext(ctx, `
				INSERT OR IGNORE INTO sprint_backlog_items (
					sprint_id, backlog_item_id, title, description, priority, status,
					provisional_story_points, story_points, acceptance_criteria, estimated_hours_hundredths
				)
				SELECT ?, backlog_item_id, title, description, priority, status,
					provisional_story_points, story_points, acceptance_criteria, estimated_hours_hundredths
				FROM sprint_backlog_items
				WHERE sprint_id = ? AND completion_verified = 0`, nextSprintID, sprintID)
			if err != nil {
				return fmt.Errorf("transfer incomplete Sprint stories: %w", err)
			}
		}

		updated, err := conn.ExecContext(ctx,
			`UPDATE sprints SET closed = 1 WHERE project_id = ? AND id = ? AND closed = 0`, projectID, sprintID)
		if err != nil {
			return fmt.Errorf("close Sprint: %w", err)
		}
		rows, err := updated.RowsAffected()
		if err != nil {
			return fmt.Errorf("count closed Sprints: %w", err)
		}
		if rows != 1 {
			return sprint.ErrAlreadyClosed
		}
		return nil
	})
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("close Sprint: %w", err)
	}
	return closed, nil
}

func (r *SprintRepository) AssignBacklogItem(ctx context.Context, projectID string, sprintID int64, value sprint.PlannedBacklogItem) error {
	return withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		var closed bool
		err := conn.QueryRowContext(ctx,
			`SELECT closed FROM sprints WHERE project_id = ? AND id = ?`, projectID, sprintID,
		).Scan(&closed)
		if errors.Is(err, sql.ErrNoRows) {
			return sprint.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("check Sprint for Backlog assignment: %w", err)
		}
		if closed {
			return sprint.ErrClosed
		}

		current, err := findBacklog(ctx, conn, projectID, value.BacklogItemID)
		if err != nil {
			return err
		}
		value.EstimatedHours = current.EstimatedHours
		criteria, err := json.Marshal(value.AcceptanceCriteria)
		if err != nil {
			return fmt.Errorf("encode planned acceptance criteria: %w", err)
		}
		result, err := conn.ExecContext(ctx, `
		INSERT OR IGNORE INTO sprint_backlog_items (
			sprint_id, backlog_item_id, title, description, priority, status,
			provisional_story_points, story_points, acceptance_criteria, estimated_hours_hundredths
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			value.SprintID, value.BacklogItemID, value.Title, value.Description, value.Priority, value.Status,
			value.ProvisionalStoryPoints, value.StoryPoints, string(criteria), value.EstimatedHours,
		)
		if err != nil {
			return fmt.Errorf("insert Sprint Backlog snapshot: %w", err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("count Sprint Backlog assignments: %w", err)
		}
		if rows == 0 {
			return sprint.ErrBacklogItemAlreadyAssigned
		}
		return nil
	})
}

func (r *SprintRepository) ListPlannedBacklogItems(ctx context.Context, projectID string, sprintID int64) ([]sprint.PlannedBacklogItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT planned.sprint_id, planned.backlog_item_id, planned.title, planned.description,
			planned.priority, planned.status, planned.completion_verified, planned.provisional_story_points,
			planned.story_points, planned.acceptance_criteria, planned.estimated_hours_hundredths
		FROM sprint_backlog_items AS planned
		JOIN sprints ON sprints.id = planned.sprint_id
		WHERE sprints.project_id = ? AND planned.sprint_id = ?
		ORDER BY planned.id`, projectID, sprintID)
	if err != nil {
		return nil, fmt.Errorf("query Sprint Backlog snapshots: %w", err)
	}
	defer rows.Close()

	plannedItems := make([]sprint.PlannedBacklogItem, 0)
	for rows.Next() {
		var value sprint.PlannedBacklogItem
		var provisional, definitive sql.NullInt64
		var completionVerified bool
		var criteria string
		if err := rows.Scan(
			&value.SprintID, &value.BacklogItemID, &value.Title, &value.Description,
			&value.Priority, &value.PlannedStatus, &completionVerified, &provisional, &definitive, &criteria, &value.EstimatedHours,
		); err != nil {
			return nil, fmt.Errorf("scan Sprint Backlog snapshot: %w", err)
		}
		value.ProvisionalStoryPoints = nullInt(provisional)
		value.StoryPoints = nullInt(definitive)
		value.Status = value.PlannedStatus
		value.AcceptanceCriteriaVerified = completionVerified
		if completionVerified {
			value.Status = string(backlog.Completed)
		}
		if err := json.Unmarshal([]byte(criteria), &value.AcceptanceCriteria); err != nil {
			return nil, fmt.Errorf("decode planned acceptance criteria: %w", err)
		}
		plannedItems = append(plannedItems, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Sprint Backlog snapshots: %w", err)
	}
	return plannedItems, nil
}

// CompleteBacklogItem updates the Backlog and its Sprint execution state atomically.
func (r *SprintRepository) CompleteBacklogItem(ctx context.Context, projectID string, sprintID int64, itemID string, verified bool) (sprint.PlannedBacklogItem, error) {
	var completed sprint.PlannedBacklogItem
	if !verified {
		return completed, backlog.ErrAcceptanceCriteriaUnverified
	}
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		var closed bool
		err := conn.QueryRowContext(ctx,
			`SELECT closed FROM sprints WHERE project_id = ? AND id = ?`, projectID, sprintID,
		).Scan(&closed)
		if errors.Is(err, sql.ErrNoRows) {
			return sprint.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("check Sprint for story completion: %w", err)
		}
		if closed {
			return sprint.ErrClosed
		}
		var exists int

		if err := conn.QueryRowContext(ctx,
			`SELECT 1 FROM backlog_items WHERE project_id = ? AND id = ?`, projectID, itemID,
		).Scan(&exists); errors.Is(err, sql.ErrNoRows) {
			return backlog.ErrNotFound
		} else if err != nil {
			return fmt.Errorf("check Backlog story for completion: %w", err)
		}

		var provisional, definitive sql.NullInt64
		var criteria string
		err = conn.QueryRowContext(ctx, `
			SELECT planned.title, planned.description, planned.priority, planned.status,
				planned.provisional_story_points, planned.story_points, planned.acceptance_criteria, planned.estimated_hours_hundredths
			FROM sprint_backlog_items AS planned
			JOIN sprints ON sprints.id = planned.sprint_id
			WHERE sprints.project_id = ? AND planned.sprint_id = ? AND planned.backlog_item_id = ?`,
			projectID, sprintID, itemID,
		).Scan(&completed.Title, &completed.Description, &completed.Priority, &completed.PlannedStatus,
			&provisional, &definitive, &criteria, &completed.EstimatedHours)
		if errors.Is(err, sql.ErrNoRows) {
			return sprint.ErrBacklogItemNotAssigned
		}
		if err != nil {
			return fmt.Errorf("find Sprint story for completion: %w", err)
		}
		if err := json.Unmarshal([]byte(criteria), &completed.AcceptanceCriteria); err != nil {
			return fmt.Errorf("decode Sprint story acceptance criteria: %w", err)
		}
		completed.SprintID = sprintID
		completed.BacklogItemID = itemID
		completed.ProvisionalStoryPoints = nullInt(provisional)
		completed.StoryPoints = nullInt(definitive)

		if _, err := conn.ExecContext(ctx,
			`UPDATE backlog_items SET status = ? WHERE project_id = ? AND id = ?`, backlog.Completed, projectID, itemID,
		); err != nil {
			return fmt.Errorf("complete Backlog story: %w", err)
		}
		if _, err := conn.ExecContext(ctx,
			`UPDATE sprint_backlog_items SET completion_verified = 1 WHERE sprint_id = ? AND backlog_item_id = ?`, sprintID, itemID,
		); err != nil {
			return fmt.Errorf("record Sprint story completion: %w", err)
		}
		return nil
	})
	if err != nil {
		return sprint.PlannedBacklogItem{}, fmt.Errorf("complete Sprint story: %w", err)
	}
	completed.Status = string(backlog.Completed)
	completed.AcceptanceCriteriaVerified = true
	return completed, nil
}
