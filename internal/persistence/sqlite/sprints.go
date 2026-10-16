package sqlite

import (
	"context"
	"database/sql"
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

func (r *SprintRepository) Create(ctx context.Context, value sprint.Sprint, updatedProject project.Project) (sprint.Sprint, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("begin create Sprint transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var currentStatus project.Status
	if err := tx.QueryRowContext(ctx,
		`SELECT status FROM projects WHERE id = ?`, value.ProjectID,
	).Scan(&currentStatus); errors.Is(err, sql.ErrNoRows) {
		return sprint.Sprint{}, project.ErrNotFound
	} else if err != nil {
		return sprint.Sprint{}, fmt.Errorf("read Sprint project: %w", err)
	}
	if currentStatus == project.Planned && updatedProject.Status != currentStatus {
		if _, err := tx.ExecContext(ctx,
			`UPDATE projects SET status = ? WHERE id = ?`, updatedProject.Status, value.ProjectID,
		); err != nil {
			return sprint.Sprint{}, fmt.Errorf("start project with first Sprint: %w", err)
		}
	}

	result, err := tx.ExecContext(ctx,
		`INSERT INTO sprints (project_id, goal) VALUES (?, ?)`, value.ProjectID, value.Goal)
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("insert Sprint: %w", err)
	}
	value.ID, err = result.LastInsertId()
	if err != nil {
		return sprint.Sprint{}, fmt.Errorf("read created Sprint ID: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return sprint.Sprint{}, fmt.Errorf("commit Sprint creation: %w", err)
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
