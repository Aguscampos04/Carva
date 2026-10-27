package integration_test

import (
	"context"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/domain/defect"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

func TestUS22SQLiteReadsProjectDefectsAndSprintAssociations(t *testing.T) {
	ctx := context.Background()
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "Defect metric repository")
	if err != nil {
		t.Fatal(err)
	}
	sprintService := sprints.NewService(sqlite.NewSprintRepository(db), nil)
	first, err := sprintService.Create(ctx, projectValue.ID, "Detection")
	if err != nil {
		t.Fatal(err)
	}
	second, err := sprintService.Create(ctx, projectValue.ID, "Resolution")
	if err != nil {
		t.Fatal(err)
	}

	insertDefect := func(id string, status defect.Status, detection, resolution any) {
		t.Helper()
		storyID := "story-" + id
		if _, err := db.ExecContext(ctx, `INSERT INTO backlog_items (
			id,project_id,title,description,priority,status,acceptance_criteria
		) VALUES (?,?,'Story','Description','Media','Pendiente','[]')`, storyID, projectValue.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := db.ExecContext(ctx, `INSERT INTO defects (
			id,project_id,backlog_item_id,description,severity,status,detection_sprint_id,resolution_sprint_id
		) VALUES (?,?,?,'Defect','Media',?,?,?)`, id, projectValue.ID, storyID, status, detection, resolution); err != nil {
			t.Fatal(err)
		}
	}
	insertDefect("split", defect.StatusResolved, first.ID, second.ID)
	insertDefect("closed", defect.StatusClosed, first.ID, first.ID)
	insertDefect("unassigned", defect.StatusOpen, nil, nil)

	data, err := sqlite.NewMetricsRepository(db).ProjectDefectData(ctx, projectValue.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Sprints) != 2 || len(data.Defects) != 3 {
		t.Fatalf("project defect data = %+v, want two Sprints and three defects", data)
	}
	if data.Sprints[0].ID != first.ID || data.Sprints[0].Goal != first.Goal || data.Sprints[0].Closed ||
		data.Sprints[1].ID != second.ID || data.Sprints[1].Goal != second.Goal {
		t.Fatalf("Sprint references = %+v", data.Sprints)
	}
	if data.Defects[0].Status != defect.StatusResolved || data.Defects[0].DetectionSprintID == nil || *data.Defects[0].DetectionSprintID != first.ID ||
		data.Defects[0].ResolutionSprintID == nil || *data.Defects[0].ResolutionSprintID != second.ID {
		t.Fatalf("defect with separate Sprints = %+v", data.Defects[0])
	}
	if data.Defects[2].DetectionSprintID != nil || data.Defects[2].ResolutionSprintID != nil {
		t.Fatalf("unassigned defect associations = %+v, want nil values", data.Defects[2])
	}
}
