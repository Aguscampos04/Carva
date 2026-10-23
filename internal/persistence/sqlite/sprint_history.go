package sqlite

import (
	"context"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

// ListClosed reads only Sprint metadata. History remains in the existing
// planning snapshots and completion records, without joining the current Backlog.
func (r *SprintRepository) ListClosed(ctx context.Context, projectID string) ([]sprint.Sprint, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, goal, closed FROM sprints WHERE project_id = ? AND closed = 1 ORDER BY id`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list closed Sprints: %w", err)
	}
	defer rows.Close()
	values := make([]sprint.Sprint, 0)
	for rows.Next() {
		var value sprint.Sprint
		if err := rows.Scan(&value.ID, &value.ProjectID, &value.Goal, &value.Closed); err != nil {
			return nil, fmt.Errorf("scan closed Sprint: %w", err)
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate closed Sprints: %w", err)
	}
	return values, nil
}
