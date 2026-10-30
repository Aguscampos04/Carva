package sqlite

import (
	"context"
	"database/sql"
	app "github.com/Aguscampos04/Carva/internal/application/reports"
	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/defect"
	"github.com/Aguscampos04/Carva/internal/domain/effort"
)

type ReportSource struct{ db *sql.DB }

func NewReportSource(db *sql.DB) *ReportSource { return &ReportSource{db} }

// Load reuses scoped repository queries. Sources may change between reads;
// this is a consultation, not an accounting closure or an atomic snapshot.
func (r *ReportSource) Load(ctx context.Context, pid string) (app.Data, error) {
	data := app.Data{Sprints: []app.SprintReport{}, Effort: []effort.Entry{}, Defects: []defect.Defect{}}
	var err error
	data.Project, err = findProject(ctx, r.db, pid)
	if err != nil {
		return app.Data{}, err
	}
	data.Members, err = NewMemberRepository(r.db).ListByProject(ctx, pid)
	if err != nil {
		return app.Data{}, err
	}
	data.Stories, err = NewBacklogRepository(r.db).ListByProject(ctx, pid)
	if err != nil {
		return app.Data{}, err
	}
	if data.Stories == nil {
		data.Stories = []backlog.Item{}
	}
	snapshots, err := NewMetricsRepository(r.db).ListProjectSprintSnapshots(ctx, pid)
	if err != nil {
		return app.Data{}, err
	}
	repo := NewSprintRepository(r.db)
	for _, value := range snapshots {
		item, err := repo.Get(ctx, pid, value.ID)
		if err != nil {
			return app.Data{}, err
		}
		stories, err := repo.ListPlannedBacklogItems(ctx, pid, value.ID)
		if err != nil {
			return app.Data{}, err
		}
		data.Sprints = append(data.Sprints, app.SprintReport{Sprint: item, Stories: stories})
	}
	for _, story := range data.Stories {
		entries, err := NewEffortRepository(r.db).List(ctx, pid, story.ID)
		if err != nil {
			return app.Data{}, err
		}
		data.Effort = append(data.Effort, entries...)
		defects, err := NewDefectRepository(r.db).ListByStory(ctx, pid, story.ID)
		if err != nil {
			return app.Data{}, err
		}
		data.Defects = append(data.Defects, defects...)
	}
	return data, nil
}

var _ app.Source = (*ReportSource)(nil)
