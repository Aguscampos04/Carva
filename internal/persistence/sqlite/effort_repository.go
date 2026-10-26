package sqlite

import (
	"context"
	"database/sql"
	"errors"
	appeffort "github.com/Aguscampos04/Carva/internal/application/effort"
	domain "github.com/Aguscampos04/Carva/internal/domain/effort"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"github.com/Aguscampos04/Carva/internal/domain/project"
)

type EffortRepository struct{ db *sql.DB }

func NewEffortRepository(db *sql.DB) *EffortRepository { return &EffortRepository{db} }

const effortColumns = `id,project_id,backlog_item_id,member_id,work_date,activity,hours_hundredths,sprint_id,deleted_at`

func scanEffort(row rowScanner) (domain.Entry, error) {
	var e domain.Entry
	err := row.Scan(&e.ID, &e.ProjectID, &e.StoryID, &e.MemberID, &e.WorkDate, &e.Activity, &e.Hours, &e.SprintID, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return e, domain.ErrNotFound
	}
	return e, err
}
func (r *EffortRepository) Get(ctx context.Context, pid, id string) (domain.Entry, error) {
	return scanEffort(r.db.QueryRowContext(ctx, `SELECT `+effortColumns+` FROM effort_entries WHERE project_id=? AND id=? AND deleted_at IS NULL`, pid, id))
}
func checkStory(ctx context.Context, reader projectReader, pid, id string) error {
	var exists bool
	if err := reader.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM backlog_items WHERE project_id=? AND id=?)`, pid, id).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return domain.ErrNotFound
	}
	return nil
}
func (r *EffortRepository) List(ctx context.Context, pid, story string) ([]domain.Entry, error) {
	if err := checkStory(ctx, r.db, pid, story); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT `+effortColumns+` FROM effort_entries WHERE project_id=? AND backlog_item_id=? AND deleted_at IS NULL ORDER BY work_date,id`, pid, story)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []domain.Entry{}
	for rows.Next() {
		e, err := scanEffort(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}
func (r *EffortRepository) Save(ctx context.Context, e domain.Entry, create bool, check func(project.Project, hours.Value) error) (domain.Entry, error) {
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		p, err := findProject(ctx, conn, e.ProjectID)
		if err != nil {
			return err
		}
		if !create {
			current, err := scanEffort(conn.QueryRowContext(ctx, `SELECT `+effortColumns+` FROM effort_entries WHERE project_id=? AND id=? AND deleted_at IS NULL`, e.ProjectID, e.ID))
			if err != nil {
				return err
			}
			if !e.HasSprint {
				e.SprintID = current.SprintID
			}
		}
		// D-30 is checked before any relationship or aggregate write.
		if p.Status == project.Finalized && p.EndDate == "" {
			return domain.ErrUnknownEnd
		}
		if err := checkStory(ctx, conn, e.ProjectID, e.StoryID); err != nil {
			return err
		}
		var exists bool
		if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM members WHERE project_id=? AND id=?)`, e.ProjectID, e.MemberID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return domain.ErrNotFound
		}
		if e.SprintID != nil {
			if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sprints WHERE project_id=? AND id=?)`, e.ProjectID, *e.SprintID).Scan(&exists); err != nil {
				return err
			}
			if !exists {
				return domain.ErrNotFound
			}
			if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sprint_backlog_items WHERE sprint_id=? AND backlog_item_id=?)`, *e.SprintID, e.StoryID).Scan(&exists); err != nil {
				return err
			}
			if !exists {
				return domain.ErrSprint
			}
		}
		var other hours.Value
		if err := conn.QueryRowContext(ctx, `SELECT COALESCE(SUM(hours_hundredths),0) FROM effort_entries WHERE project_id=? AND member_id=? AND work_date=? AND deleted_at IS NULL AND id<>?`, e.ProjectID, e.MemberID, e.WorkDate, e.ID).Scan(&other); err != nil {
			return err
		}
		if err := check(p, other); err != nil {
			return err
		}
		if create {
			_, err = conn.ExecContext(ctx, `INSERT INTO effort_entries(id,project_id,backlog_item_id,member_id,work_date,activity,hours_hundredths,sprint_id) VALUES(?,?,?,?,?,?,?,?)`, e.ID, e.ProjectID, e.StoryID, e.MemberID, e.WorkDate, e.Activity, e.Hours, e.SprintID)
		} else {
			_, err = conn.ExecContext(ctx, `UPDATE effort_entries SET backlog_item_id=?,member_id=?,work_date=?,activity=?,hours_hundredths=?,sprint_id=? WHERE project_id=? AND id=? AND deleted_at IS NULL`, e.StoryID, e.MemberID, e.WorkDate, e.Activity, e.Hours, e.SprintID, e.ProjectID, e.ID)
		}
		return err
	})
	return e, err
}
func (r *EffortRepository) Delete(ctx context.Context, pid, id, stamp string) error {
	return withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		var exists bool
		if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM effort_entries WHERE project_id=? AND id=?)`, pid, id).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return domain.ErrNotFound
		}
		_, err := conn.ExecContext(ctx, `UPDATE effort_entries SET deleted_at=? WHERE project_id=? AND id=? AND deleted_at IS NULL`, stamp, pid, id)
		return err
	})
}
func (r *EffortRepository) Sprints(ctx context.Context, pid, story string) ([]appeffort.SprintOption, error) {
	if err := checkStory(ctx, r.db, pid, story); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT s.id,s.closed FROM sprints s JOIN sprint_backlog_items b ON b.sprint_id=s.id WHERE s.project_id=? AND b.backlog_item_id=? ORDER BY s.id`, pid, story)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []appeffort.SprintOption{}
	for rows.Next() {
		var v appeffort.SprintOption
		if err := rows.Scan(&v.ID, &v.Closed); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}
