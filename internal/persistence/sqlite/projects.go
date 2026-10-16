package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/project"
)

type ProjectRepository struct{ db *sql.DB }

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, value project.Project) (project.Project, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO projects (name, status) VALUES (?, ?)`, value.Name, value.Status)
	if err != nil {
		return project.Project{}, fmt.Errorf("insert project: %w", err)
	}
	value.ID, err = result.LastInsertId()
	if err != nil {
		return project.Project{}, fmt.Errorf("read created project ID: %w", err)
	}
	return value, nil
}

func (r *ProjectRepository) Get(ctx context.Context, id int64) (project.Project, error) {
	var value project.Project
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, status FROM projects WHERE id = ?`, id,
	).Scan(&value.ID, &value.Name, &value.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return project.Project{}, project.ErrNotFound
	}
	if err != nil {
		return project.Project{}, fmt.Errorf("get project: %w", err)
	}
	return value, nil
}
