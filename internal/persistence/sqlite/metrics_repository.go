package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Aguscampos04/Carva/internal/domain/metrics"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type MetricsRepository struct{ db *sql.DB }

func NewMetricsRepository(db *sql.DB) *MetricsRepository { return &MetricsRepository{db: db} }

func (r *MetricsRepository) ListProjectSprintSnapshots(ctx context.Context, projectID string) ([]domain.SprintSnapshot, error) {
	if err := projectExists(ctx, r.db, projectID); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT sprints.id, sprints.goal, sprints.closed,
			planned.id, planned.story_points, planned.completion_verified
		FROM sprints
	LEFT JOIN sprint_backlog_items AS planned ON planned.sprint_id = sprints.id
		WHERE sprints.project_id = ?
		ORDER BY sprints.id, planned.id`, projectID)
	if err != nil {
		return nil, fmt.Errorf("query Sprint snapshots for metrics: %w", err)
	}
	defer rows.Close()

	var result []domain.SprintSnapshot
	for rows.Next() {
		var sprintID int64
		var goal string
		var closed bool
		var plannedID, storyPoints sql.NullInt64
		var completed sql.NullBool
		if err := rows.Scan(&sprintID, &goal, &closed, &plannedID, &storyPoints, &completed); err != nil {
			return nil, fmt.Errorf("scan Sprint snapshot for metrics: %w", err)
		}
		if len(result) == 0 || result[len(result)-1].ID != sprintID {
			result = append(result, domain.SprintSnapshot{ID: sprintID, Goal: goal, Closed: closed, Items: make([]domain.StorySnapshot, 0)})
		}
		if plannedID.Valid {
			item := domain.StorySnapshot{Completed: completed.Valid && completed.Bool}
			if storyPoints.Valid {
				value := int(storyPoints.Int64)
				item.StoryPoints = &value
			}
			last := len(result) - 1
			result[last].Items = append(result[last].Items, item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Sprint snapshots for metrics: %w", err)
	}
	return result, nil
}

func (r *MetricsRepository) GetSprintSnapshot(ctx context.Context, projectID string, sprintID int64) (domain.SprintSnapshot, error) {
	if sprintID < 1 {
		return domain.SprintSnapshot{}, sprint.ErrNotFound
	}
	snapshots, err := r.ListProjectSprintSnapshots(ctx, projectID)
	if err != nil {
		return domain.SprintSnapshot{}, err
	}
	for _, snapshot := range snapshots {
		if snapshot.ID == sprintID {
			return snapshot, nil
		}
	}
	return domain.SprintSnapshot{}, sprint.ErrNotFound
}

func projectExists(ctx context.Context, db *sql.DB, projectID string) error {
	var found int
	err := db.QueryRowContext(ctx, `SELECT 1 FROM projects WHERE id = ?`, projectID).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return project.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("find project for metrics: %w", err)
	}
	return nil
}
