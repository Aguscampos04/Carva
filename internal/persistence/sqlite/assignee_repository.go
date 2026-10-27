package sqlite

import (
	"context"
	"database/sql"
	"errors"

	app "github.com/Aguscampos04/Carva/internal/application/backlog"
	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
)

type AssigneeRepository struct{ db *sql.DB }

func NewAssigneeRepository(db *sql.DB) *AssigneeRepository { return &AssigneeRepository{db: db} }
func (r *AssigneeRepository) Get(ctx context.Context, pid, id string) (app.Assignee, error) {
	var member sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT a.member_id FROM backlog_items b LEFT JOIN story_assignees a ON a.backlog_item_id=b.id WHERE b.id=? AND b.project_id=?`, id, pid).Scan(&member)
	if errors.Is(err, sql.ErrNoRows) {
		return app.Assignee{}, domain.ErrNotFound
	}
	if err != nil {
		return app.Assignee{}, err
	}
	if member.Valid {
		return app.Assignee{MemberID: &member.String}, nil
	}
	return app.Assignee{}, nil
}
func (r *AssigneeRepository) Set(ctx context.Context, pid, id string, memberID *string) error {
	return withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		var exists bool
		if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM backlog_items WHERE id=? AND project_id=?)`, id, pid).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return domain.ErrNotFound
		}
		if memberID == nil {
			_, err := conn.ExecContext(ctx, `DELETE FROM story_assignees WHERE backlog_item_id=?`, id)
			return err
		}
		if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM members WHERE id=? AND project_id=?)`, *memberID, pid).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return app.ErrInvalidAssignee
		}
		_, err := conn.ExecContext(ctx, `INSERT INTO story_assignees(backlog_item_id,member_id) VALUES(?,?) ON CONFLICT(backlog_item_id) DO UPDATE SET member_id=excluded.member_id`, id, *memberID)
		return err
	})
}
