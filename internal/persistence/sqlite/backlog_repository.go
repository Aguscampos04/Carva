package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Aguscampos04/Carva/internal/domain/project"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
)

// BacklogRepository stores Product Backlog items in SQLite.
type BacklogRepository struct {
	db *sql.DB
}

// NewBacklogRepository constructs a SQLite Backlog repository.
func NewBacklogRepository(db *sql.DB) *BacklogRepository {
	return &BacklogRepository{db: db}
}

// ProjectExists reports whether projectID identifies a persisted project.
func (r *BacklogRepository) ProjectExists(ctx context.Context, projectID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)`, projectID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("check project existence: %w", err)
	}
	return exists, nil
}

// Create inserts a Backlog item.
func (r *BacklogRepository) Create(ctx context.Context, item domain.Item) error {
	criteria, err := json.Marshal(item.AcceptanceCriteria)
	if err != nil {
		return fmt.Errorf("encode acceptance criteria: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO backlog_items (
			id, project_id, title, description, priority, status,
			provisional_story_points, story_points, acceptance_criteria, estimated_hours_hundredths
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.ProjectID, item.Title, item.Description, item.Priority, item.Status,
		item.ProvisionalStoryPoints, item.StoryPoints, string(criteria), item.EstimatedHours,
	)
	if err != nil {
		return fmt.Errorf("insert backlog item: %w", err)
	}
	return nil
}

// Reads use a correlated aggregate, never a join that duplicates effort across Sprints.
const backlogSelect = `SELECT b.id,b.project_id,b.title,b.description,b.priority,b.status,
 b.provisional_story_points,b.story_points,b.acceptance_criteria,b.estimated_hours_hundredths,
 COALESCE((SELECT SUM(e.hours_hundredths) FROM effort_entries e WHERE e.project_id=b.project_id AND e.backlog_item_id=b.id AND e.deleted_at IS NULL),0)
 FROM backlog_items b `

func findBacklog(ctx context.Context, q projectReader, pid, id string) (domain.Item, error) {
	value, err := scanBacklogItem(q.QueryRowContext(ctx, backlogSelect+`WHERE b.project_id=? AND b.id=?`, pid, id))
	if errors.Is(err, sql.ErrNoRows) {
		return value, domain.ErrNotFound
	}
	return value, err
}
func (r *BacklogRepository) FindByID(ctx context.Context, pid, id string) (domain.Item, error) {
	return findBacklog(ctx, r.db, pid, id)
}
func (r *BacklogRepository) Mutate(ctx context.Context, pid, id string, change func(domain.Item) (domain.Item, error)) (domain.Item, error) {
	var next domain.Item
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		if _, err := findProject(ctx, conn, pid); errors.Is(err, project.ErrNotFound) {
			return domain.ErrProjectNotFound
		} else if err != nil {
			return err
		}
		current, err := findBacklog(ctx, conn, pid, id)
		if err != nil {
			return err
		}
		next, err = change(current)
		if err != nil {
			return err
		}
		criteria, err := json.Marshal(next.AcceptanceCriteria)
		if err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, `UPDATE backlog_items SET title=?,description=?,priority=?,provisional_story_points=?,acceptance_criteria=?,estimated_hours_hundredths=? WHERE project_id=? AND id=?`, next.Title, next.Description, next.Priority, next.ProvisionalStoryPoints, string(criteria), next.EstimatedHours, pid, id)
		if err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, `UPDATE sprint_backlog_items SET estimated_hours_hundredths=? WHERE backlog_item_id=? AND sprint_id IN (SELECT id FROM sprints WHERE project_id=?)`, next.EstimatedHours, id, pid)
		return err
	})
	return next, err
}
func (r *BacklogRepository) ListByProject(ctx context.Context, pid string) ([]domain.Item, error) {
	rows, err := r.db.QueryContext(ctx, backlogSelect+`WHERE b.project_id=?`, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []domain.Item{}
	for rows.Next() {
		v, err := scanBacklogItem(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

type rowScanner interface {
	Scan(...any) error
}

func scanBacklogItem(scanner rowScanner) (domain.Item, error) {
	var item domain.Item
	var priority, status, criteria string
	var provisional, definitive sql.NullInt64
	if err := scanner.Scan(
		&item.ID, &item.ProjectID, &item.Title, &item.Description, &priority, &status,
		&provisional, &definitive, &criteria, &item.EstimatedHours, &item.ActualHours,
	); err != nil {
		return domain.Item{}, err
	}
	item.Priority = domain.Priority(priority)
	item.Status = domain.Status(status)
	item.ProvisionalStoryPoints = nullInt(provisional)
	item.StoryPoints = nullInt(definitive)
	if err := json.Unmarshal([]byte(criteria), &item.AcceptanceCriteria); err != nil {
		return domain.Item{}, fmt.Errorf("decode acceptance criteria: %w", err)
	}
	return item, nil
}

func nullInt(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	converted := int(value.Int64)
	return &converted
}

var _ appbacklog.Repository = (*BacklogRepository)(nil)
