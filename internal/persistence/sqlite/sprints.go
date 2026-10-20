package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type SprintRepository struct{ db *sql.DB }

func NewSprintRepository(db *sql.DB) *SprintRepository {
	return &SprintRepository{db: db}
}

func (r *SprintRepository) GetProject(ctx context.Context, id string) (project.Project, error) {
	var value project.Project
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, status FROM projects WHERE id = ?`, id,
	).Scan(&value.ID, &value.Name, &value.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return project.Project{}, project.ErrNotFound
	}
	if err != nil {
		return project.Project{}, fmt.Errorf("get Sprint project: %w", err)
	}
	return value, nil
}

func (r *SprintRepository) Create(ctx context.Context, value sprint.Sprint, transition func(project.Project) (project.Project, error)) (sprint.Sprint, error) {
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		current, err := readProject(conn.QueryRowContext(ctx, `SELECT id,name,status,start_date,end_date FROM projects WHERE id=?`, value.ProjectID))
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
		`SELECT id, project_id, goal FROM sprints WHERE project_id = ? AND id = ?`, projectID, sprintID,
	).Scan(&value.ID, &value.ProjectID, &value.Goal)
	if errors.Is(err, sql.ErrNoRows) {
		return sprint.Sprint{}, sprint.ErrNotFound
	}
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("get Sprint: %w", err)
	}
	return value, nil
}

func (r *SprintRepository) AssignBacklogItem(ctx context.Context, projectID string, sprintID int64, value sprint.PlannedBacklogItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin Sprint Backlog assignment: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var exists int
	err = tx.QueryRowContext(ctx,
		`SELECT 1 FROM sprints WHERE project_id = ? AND id = ?`, projectID, sprintID,
	).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return sprint.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("check Sprint for Backlog assignment: %w", err)
	}

	criteria, err := json.Marshal(value.AcceptanceCriteria)
	if err != nil {
		return fmt.Errorf("encode planned acceptance criteria: %w", err)
	}
	result, err := tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO sprint_backlog_items (
			sprint_id, backlog_item_id, title, description, priority, status,
			provisional_story_points, story_points, acceptance_criteria
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		value.SprintID, value.BacklogItemID, value.Title, value.Description, value.Priority, value.Status,
		value.ProvisionalStoryPoints, value.StoryPoints, string(criteria),
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
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit Sprint Backlog assignment: %w", err)
	}
	return nil
}

func (r *SprintRepository) ListPlannedBacklogItems(ctx context.Context, projectID string, sprintID int64) ([]sprint.PlannedBacklogItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT planned.sprint_id, planned.backlog_item_id, planned.title, planned.description,
			planned.priority, planned.status, planned.provisional_story_points,
			planned.story_points, planned.acceptance_criteria
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
		var criteria string
		if err := rows.Scan(
			&value.SprintID, &value.BacklogItemID, &value.Title, &value.Description,
			&value.Priority, &value.Status, &provisional, &definitive, &criteria,
		); err != nil {
			return nil, fmt.Errorf("scan Sprint Backlog snapshot: %w", err)
		}
		value.ProvisionalStoryPoints = nullInt(provisional)
		value.StoryPoints = nullInt(definitive)
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
