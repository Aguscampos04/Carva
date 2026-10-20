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
	return findProject(ctx, r.db, id)
}

// Both a pool and a reserved transactional connection support this read port.
type projectReader interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func findProject(ctx context.Context, reader projectReader, id string) (project.Project, error) {
	return readProject(reader.QueryRowContext(ctx, `SELECT id,name,status,start_date,end_date FROM projects WHERE id=?`, id))
}

func readProject(row *sql.Row) (project.Project, error) {
	var found project.Project
	var status string
	var start, end sql.NullString
	err := row.Scan(&found.ID, &found.Name, &status, &start, &end)
	if errors.Is(err, sql.ErrNoRows) {
		return project.Project{}, project.ErrNotFound
	}
	if err != nil {
		return project.Project{}, fmt.Errorf("select project: %w", err)
	}
	found.Status = project.Status(status)
	if start.Valid {
		found.StartDate, err = project.ParseDate(start.String)
		if err != nil {
			return project.Project{}, err
		}
	}
	if end.Valid {
		found.EndDate, err = project.ParseDate(end.String)
		if err != nil {
			return project.Project{}, err
		}
	}
	return found, nil
}

// Mutate reads and validates the latest state under the same write transaction.
func (r *ProjectRepository) Mutate(ctx context.Context, id string, change func(project.Project) (project.Project, error)) (project.Project, error) {
	var next project.Project
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		current, err := findProject(ctx, conn, id)
		if err != nil {
			return err
		}
		next, err = change(current)
		if err != nil {
			return err
		}
		if next == current {
			return nil
		}
		_, err = conn.ExecContext(ctx, `UPDATE projects SET name=?,status=?,start_date=?,end_date=? WHERE id=?`, next.Name, next.Status, nullableDate(next.StartDate), nullableDate(next.EndDate), id)
		return err
	})
	if err != nil {
		return project.Project{}, fmt.Errorf("mutate project: %w", err)
	}
	return next, nil
}

func nullableDate(date project.Date) any {
	if date == "" {
		return nil
	}
	return string(date)
}
