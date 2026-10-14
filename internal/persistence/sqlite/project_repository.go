package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/project"
)

// ProjectRepository stores project entities in SQLite.
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository constructs a SQLite project repository.
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Save inserts a project without changing an existing project with the same ID.
func (r *ProjectRepository) Save(ctx context.Context, p project.Project) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO projects (id, name, status) VALUES (?, ?, ?)`, p.ID, p.Name, p.Status,
	)
	if err != nil {
		return fmt.Errorf("insert project: %w", err)
	}
	return nil
}

// FindByID retrieves a project by its stable identifier.
func (r *ProjectRepository) FindByID(ctx context.Context, id string) (project.Project, error) {
	var found project.Project
	var status string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, status FROM projects WHERE id = ?`, id,
	).Scan(&found.ID, &found.Name, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return project.Project{}, project.ErrNotFound
	}
	if err != nil {
		return project.Project{}, fmt.Errorf("select project: %w", err)
	}
	found.Status = project.Status(status)
	return found, nil
}
