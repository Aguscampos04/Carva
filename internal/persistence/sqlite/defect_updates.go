package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Aguscampos04/Carva/internal/domain/defect"
)

// Update reads the latest scoped row and merges the patch under the same write
// transaction. A retry re-reads current fields and never replays an external effect.
func (r *DefectRepository) Update(ctx context.Context, projectID, storyID, id string, change domain.Patch) (domain.Defect, error) {
	if err := change.Validate(); err != nil {
		return domain.Defect{}, err
	}
	var updated domain.Defect
	err := withImmediate(ctx, r.db, func(conn *sql.Conn) error {
		current, err := scanDefect(conn.QueryRowContext(ctx, `SELECT `+defectColumns+` FROM defects WHERE project_id=? AND backlog_item_id=? AND id=?`, projectID, storyID, id))
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("find defect to update: %w", err)
		}
		updated, err = current.Update(change)
		if err != nil {
			return err
		}
		for _, sprintID := range []*int64{updated.DetectionSprintID, updated.ResolutionSprintID} {
			if sprintID == nil {
				continue
			}
			var exists bool
			if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sprints WHERE project_id=? AND id=?)`, projectID, *sprintID).Scan(&exists); err != nil {
				return fmt.Errorf("find updated defect Sprint: %w", err)
			}
			if !exists {
				return domain.ErrSprintNotFound
			}
		}
		_, err = conn.ExecContext(ctx, `UPDATE defects SET status=?,detection_sprint_id=?,resolution_sprint_id=? WHERE project_id=? AND backlog_item_id=? AND id=?`, updated.Status, nullableSprint(updated.DetectionSprintID), nullableSprint(updated.ResolutionSprintID), projectID, storyID, id)
		return err
	})
	if err != nil {
		return domain.Defect{}, fmt.Errorf("persist defect update: %w", err)
	}
	return updated, nil
}
