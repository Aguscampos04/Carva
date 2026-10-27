package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/hours"
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

func (r *MetricsRepository) ProjectDefectData(ctx context.Context, projectID string) (domain.ProjectDefectData, error) {
	if err := projectExists(ctx, r.db, projectID); err != nil {
		return domain.ProjectDefectData{}, err
	}
	result := domain.ProjectDefectData{
		Sprints: make([]domain.DefectSprint, 0),
		Defects: make([]domain.DefectSnapshot, 0),
	}
	sprintRows, err := r.db.QueryContext(ctx, `SELECT id, goal, closed FROM sprints WHERE project_id=? ORDER BY id`, projectID)
	if err != nil {
		return domain.ProjectDefectData{}, fmt.Errorf("query Sprints for defect metrics: %w", err)
	}
	for sprintRows.Next() {
		var value domain.DefectSprint
		if err := sprintRows.Scan(&value.ID, &value.Goal, &value.Closed); err != nil {
			sprintRows.Close()
			return domain.ProjectDefectData{}, fmt.Errorf("scan Sprint for defect metrics: %w", err)
		}
		result.Sprints = append(result.Sprints, value)
	}
	if err := sprintRows.Err(); err != nil {
		sprintRows.Close()
		return domain.ProjectDefectData{}, fmt.Errorf("iterate Sprints for defect metrics: %w", err)
	}
	if err := sprintRows.Close(); err != nil {
		return domain.ProjectDefectData{}, fmt.Errorf("close Sprint query for defect metrics: %w", err)
	}

	defectRows, err := r.db.QueryContext(ctx, `
		SELECT status, detection_sprint_id, resolution_sprint_id
		FROM defects WHERE project_id=? ORDER BY rowid`, projectID)
	if err != nil {
		return domain.ProjectDefectData{}, fmt.Errorf("query defects for metrics: %w", err)
	}
	defer defectRows.Close()
	for defectRows.Next() {
		var value domain.DefectSnapshot
		var detection, resolution sql.NullInt64
		if err := defectRows.Scan(&value.Status, &detection, &resolution); err != nil {
			return domain.ProjectDefectData{}, fmt.Errorf("scan defect for metrics: %w", err)
		}
		value.DetectionSprintID = nullableInt64(detection)
		value.ResolutionSprintID = nullableInt64(resolution)
		result.Defects = append(result.Defects, value)
	}
	if err := defectRows.Err(); err != nil {
		return domain.ProjectDefectData{}, fmt.Errorf("iterate defects for metrics: %w", err)
	}
	return result, nil
}

func (r *MetricsRepository) ProjectHoursData(ctx context.Context, projectID string) (domain.ProjectHoursData, error) {
	if err := projectExists(ctx, r.db, projectID); err != nil {
		return domain.ProjectHoursData{}, err
	}
	result := domain.ProjectHoursData{Sprints: make([]domain.SprintHoursData, 0)}

	storyRows, err := r.db.QueryContext(ctx, `SELECT estimated_hours_hundredths FROM backlog_items WHERE project_id=? ORDER BY id`, projectID)
	if err != nil {
		return domain.ProjectHoursData{}, fmt.Errorf("query backlog estimates for hour metrics: %w", err)
	}
	for storyRows.Next() {
		var estimate sql.NullInt64
		if err := storyRows.Scan(&estimate); err != nil {
			storyRows.Close()
			return domain.ProjectHoursData{}, fmt.Errorf("scan backlog estimate for hour metrics: %w", err)
		}
		result.Estimates = append(result.Estimates, nullableHours(estimate))
	}
	if err := storyRows.Err(); err != nil {
		storyRows.Close()
		return domain.ProjectHoursData{}, fmt.Errorf("iterate backlog estimates for hour metrics: %w", err)
	}
	if err := storyRows.Close(); err != nil {
		return domain.ProjectHoursData{}, fmt.Errorf("close backlog estimate query for hour metrics: %w", err)
	}

	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(hours_hundredths),0) FROM effort_entries WHERE project_id=? AND deleted_at IS NULL`, projectID).Scan(&result.ActualHours); err != nil {
		return domain.ProjectHoursData{}, fmt.Errorf("sum project effort for hour metrics: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.goal, s.closed, planned.backlog_item_id, planned.estimated_hours_hundredths,
			COALESCE((SELECT SUM(e.hours_hundredths) FROM effort_entries e WHERE e.project_id=s.project_id AND e.sprint_id=s.id AND e.deleted_at IS NULL),0)
		FROM sprints s
		LEFT JOIN sprint_backlog_items AS planned ON planned.sprint_id=s.id
		WHERE s.project_id=?
		ORDER BY s.id, planned.backlog_item_id`, projectID)
	if err != nil {
		return domain.ProjectHoursData{}, fmt.Errorf("query Sprint snapshots for hour metrics: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sprintID int64
		var goal string
		var closed bool
		var storyID sql.NullString
		var estimate sql.NullInt64
		var actual int64
		if err := rows.Scan(&sprintID, &goal, &closed, &storyID, &estimate, &actual); err != nil {
			return domain.ProjectHoursData{}, fmt.Errorf("scan Sprint snapshot for hour metrics: %w", err)
		}
		if len(result.Sprints) == 0 || result.Sprints[len(result.Sprints)-1].ID != sprintID {
			result.Sprints = append(result.Sprints, domain.SprintHoursData{ID: sprintID, Goal: goal, Closed: closed, Estimates: make([]*hours.Value, 0)})
			result.Sprints[len(result.Sprints)-1].ActualHours = hours.Value(actual)
		}
		if storyID.Valid {
			last := len(result.Sprints) - 1
			result.Sprints[last].Estimates = append(result.Sprints[last].Estimates, nullableHours(estimate))
		}
	}
	if err := rows.Err(); err != nil {
		return domain.ProjectHoursData{}, fmt.Errorf("iterate Sprint snapshots for hour metrics: %w", err)
	}
	return result, nil
}

func nullableHours(value sql.NullInt64) *hours.Value {
	if !value.Valid {
		return nil
	}
	result := hours.Value(value.Int64)
	return &result
}

func nullableInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
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
