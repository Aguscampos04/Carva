package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

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
			provisional_story_points, story_points, acceptance_criteria
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.ProjectID, item.Title, item.Description, item.Priority, item.Status,
		item.ProvisionalStoryPoints, item.StoryPoints, string(criteria),
	)
	if err != nil {
		return fmt.Errorf("insert backlog item: %w", err)
	}
	return nil
}

// FindByID retrieves an item that belongs to projectID.
func (r *BacklogRepository) FindByID(ctx context.Context, projectID, id string) (domain.Item, error) {
	item, err := scanBacklogItem(r.db.QueryRowContext(ctx, `
		SELECT id, project_id, title, description, priority, status,
			provisional_story_points, story_points, acceptance_criteria
		FROM backlog_items WHERE project_id = ? AND id = ?`, projectID, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Item{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Item{}, fmt.Errorf("select backlog item: %w", err)
	}
	return item, nil
}

// Update replaces one item's editable and preserved fields.
func (r *BacklogRepository) Update(ctx context.Context, item domain.Item) error {
	criteria, err := json.Marshal(item.AcceptanceCriteria)
	if err != nil {
		return fmt.Errorf("encode acceptance criteria: %w", err)
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE backlog_items SET title = ?, description = ?, priority = ?, status = ?,
			provisional_story_points = ?, story_points = ?, acceptance_criteria = ?
		WHERE project_id = ? AND id = ?`,
		item.Title, item.Description, item.Priority, item.Status,
		item.ProvisionalStoryPoints, item.StoryPoints, string(criteria), item.ProjectID, item.ID,
	)
	if err != nil {
		return fmt.Errorf("update backlog item: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count updated backlog items: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ListByProject returns every Backlog item for projectID.
func (r *BacklogRepository) ListByProject(ctx context.Context, projectID string) ([]domain.Item, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, project_id, title, description, priority, status,
			provisional_story_points, story_points, acceptance_criteria
		FROM backlog_items WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("query backlog items: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Item, 0)
	for rows.Next() {
		item, err := scanBacklogItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan backlog item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate backlog items: %w", err)
	}
	return items, nil
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
		&provisional, &definitive, &criteria,
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
