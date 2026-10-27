package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	app "github.com/Aguscampos04/Carva/internal/application/defects"
	domain "github.com/Aguscampos04/Carva/internal/domain/defect"
	"github.com/Aguscampos04/Carva/internal/domain/project"
)

type DefectRepository struct{ db *sql.DB }

func NewDefectRepository(db *sql.DB) *DefectRepository { return &DefectRepository{db: db} }

const defectColumns = `id,project_id,backlog_item_id,description,severity,status,detection_sprint_id,resolution_sprint_id`

func scanDefect(row rowScanner) (domain.Defect, error) {
	var item domain.Defect
	var detection, resolution sql.NullInt64
	err := row.Scan(&item.ID, &item.ProjectID, &item.BacklogItemID, &item.Description, &item.Severity, &item.Status, &detection, &resolution)
	if errors.Is(err, project.ErrNotFound) {
		return item, domain.ErrBacklogItemNotFound
	}
	if err != nil {
		return item, err
	}
	if detection.Valid {
		item.DetectionSprintID = &detection.Int64
	}
	if resolution.Valid {
		item.ResolutionSprintID = &resolution.Int64
	}
	return item, nil
}

func (r *DefectRepository) Create(ctx context.Context, item domain.Defect) error {
	return withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		if _, err := findProject(ctx, conn, item.ProjectID); errors.Is(err, project.ErrNotFound) {
			return domain.ErrProjectNotFound
		} else if err != nil {
			return fmt.Errorf("find defect project: %w", err)
		}
		var exists bool
		if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM backlog_items WHERE project_id=? AND id=?)`, item.ProjectID, item.BacklogItemID).Scan(&exists); err != nil {
			return fmt.Errorf("find defect story: %w", err)
		}
		if !exists {
			return domain.ErrBacklogItemNotFound
		}
		if err := validateDefectSprints(ctx, conn, item); err != nil {
			return err
		}
		_, err := conn.ExecContext(ctx, `INSERT INTO defects(id,project_id,backlog_item_id,description,severity,status,detection_sprint_id,resolution_sprint_id) VALUES(?,?,?,?,?,?,?,?)`,
			item.ID, item.ProjectID, item.BacklogItemID, item.Description, item.Severity, item.Status,
			nullableSprint(item.DetectionSprintID), nullableSprint(item.ResolutionSprintID))
		return err
	})
}

func nullableSprint(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}

func (r *DefectRepository) ListByStory(ctx context.Context, projectID, backlogItemID string) ([]domain.Defect, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM backlog_items WHERE project_id=? AND id=?)`, projectID, backlogItemID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("find defect story: %w", err)
	}
	if !exists {
		var projectExists bool
		if err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM projects WHERE id=?)`, projectID).Scan(&projectExists); err != nil {
			return nil, fmt.Errorf("find defect project: %w", err)
		}
		if !projectExists {
			return nil, domain.ErrProjectNotFound
		}
		return nil, domain.ErrBacklogItemNotFound
	}
	rows, err := r.db.QueryContext(ctx, `SELECT `+defectColumns+` FROM defects WHERE project_id=? AND backlog_item_id=? ORDER BY rowid`, projectID, backlogItemID)
	if err != nil {
		return nil, fmt.Errorf("query defects: %w", err)
	}
	defer rows.Close()
	items := make([]domain.Defect, 0)
	for rows.Next() {
		item, err := scanDefect(rows)
		if err != nil {
			return nil, fmt.Errorf("scan defect: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate defects: %w", err)
	}
	return items, nil
}

var _ app.Repository = (*DefectRepository)(nil)
