package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/application/members"
	"github.com/Aguscampos04/Carva/internal/domain/member"
)

// MemberRepository persists project-scoped identities without name deduplication.
type MemberRepository struct{ db *sql.DB }

func NewMemberRepository(db *sql.DB) *MemberRepository { return &MemberRepository{db: db} }

func (r *MemberRepository) ProjectExists(ctx context.Context, projectID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)`, projectID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("find member project: %w", err)
	}
	return exists, nil
}

func (r *MemberRepository) Create(ctx context.Context, m member.Member) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO members(id,project_id,name) VALUES(?,?,?)`, m.ID, m.ProjectID, m.Name)
	if err != nil {
		return fmt.Errorf("insert member: %w", err)
	}
	return nil
}

func (r *MemberRepository) ListByProject(ctx context.Context, projectID string) ([]member.Member, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,project_id,name FROM members WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("query members: %w", err)
	}
	defer rows.Close()
	items := make([]member.Member, 0)
	for rows.Next() {
		var m member.Member
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.Name); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate members: %w", err)
	}
	return items, nil
}

var _ members.Repository = (*MemberRepository)(nil)
