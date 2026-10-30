// Package reports composes read-only reports from the existing sources.
package reports

import (
	"context"
	"errors"
	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/defect"
	"github.com/Aguscampos04/Carva/internal/domain/effort"
	"github.com/Aguscampos04/Carva/internal/domain/member"
	"github.com/Aguscampos04/Carva/internal/domain/metrics"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

var ErrInvalidSprintID = errors.New("Sprint ID must be a positive integer")
var ErrInconsistentSources = errors.New("report sources changed; retry consultation")

type Summary struct {
	Points  metrics.StoryPointMetrics `json:"story_points"`
	Hours   metrics.HoursMetrics      `json:"hours"`
	Defects metrics.DefectMetrics     `json:"defects"`
}
type SprintReport struct {
	Sprint  sprint.Sprint               `json:"sprint"`
	Stories []sprint.PlannedBacklogItem `json:"stories"`
	Metrics Summary                     `json:"metrics"`
}
type Data struct {
	Project project.Project `json:"project"`
	Members []member.Member `json:"members"`
	Stories []backlog.Item  `json:"stories"`
	Sprints []SprintReport  `json:"sprints"`
	Effort  []effort.Entry  `json:"effort"`
	Defects []defect.Defect `json:"defects"`
}
type Report struct {
	Scope   string          `json:"scope"`
	Project project.Project `json:"project"`
	Members []member.Member `json:"members"`
	Stories []backlog.Item  `json:"stories"`
	Sprints []SprintReport  `json:"sprints"`
	Effort  []effort.Entry  `json:"effort"`
	Defects []defect.Defect `json:"defects"`
	Metrics Summary         `json:"metrics"`
}
type Source interface {
	Load(context.Context, string) (Data, error)
}
type Metrics interface {
	ProjectStoryPoints(context.Context, string) (metrics.ProjectStoryPointMetrics, error)
	ProjectHours(context.Context, string) (metrics.ProjectHoursMetrics, error)
	ProjectDefects(context.Context, string) (metrics.ProjectDefectMetrics, error)
}
type Service struct {
	source  Source
	metrics Metrics
}

func NewService(source Source, metrics Metrics) *Service { return &Service{source, metrics} }
func (s *Service) Project(ctx context.Context, pid string) (Report, error) {
	data, err := s.source.Load(ctx, pid)
	if err != nil {
		return Report{}, err
	}
	points, err := s.metrics.ProjectStoryPoints(ctx, pid)
	if err != nil {
		return Report{}, err
	}
	hrs, err := s.metrics.ProjectHours(ctx, pid)
	if err != nil {
		return Report{}, err
	}
	defects, err := s.metrics.ProjectDefects(ctx, pid)
	if err != nil {
		return Report{}, err
	}
	result := Report{Scope: "project", Project: data.Project, Members: data.Members, Stories: data.Stories, Sprints: data.Sprints, Effort: data.Effort, Defects: data.Defects, Metrics: Summary{points.Metrics, hrs.Metrics, defects.Metrics}}
	for i := range result.Sprints {
		id := result.Sprints[i].Sprint.ID
		found := 0
		for _, v := range points.Sprints {
			if v.ID == id {
				result.Sprints[i].Metrics.Points = v.Metrics
				found++
			}
		}
		for _, v := range hrs.Sprints {
			if v.ID == id {
				result.Sprints[i].Metrics.Hours = v.Metrics
				found++
			}
		}
		for _, v := range defects.Sprints {
			if v.ID == id {
				result.Sprints[i].Metrics.Defects = v.Metrics
				found++
			}
		}
		if found != 3 {
			return Report{}, ErrInconsistentSources
		}
	}
	return result, nil
}
func (s *Service) Sprint(ctx context.Context, pid string, id int64) (Report, error) {
	if id < 1 {
		return Report{}, ErrInvalidSprintID
	}
	report, err := s.Project(ctx, pid)
	if err != nil {
		return Report{}, err
	}
	return ForSprint(report, id)
}

// ForSprint narrows scope without rebuilding historical planning from Backlog.
func ForSprint(source Report, id int64) (Report, error) {
	var selected *SprintReport
	for i := range source.Sprints {
		if source.Sprints[i].Sprint.ID == id {
			selected = &source.Sprints[i]
			break
		}
	}
	if selected == nil {
		return Report{}, sprint.ErrNotFound
	}
	result := source
	result.Scope = "sprint"
	result.Stories = []backlog.Item{}
	result.Sprints = []SprintReport{*selected}
	result.Metrics = selected.Metrics
	result.Effort = []effort.Entry{}
	result.Defects = []defect.Defect{}
	for _, v := range source.Effort {
		if v.SprintID != nil && *v.SprintID == id {
			result.Effort = append(result.Effort, v)
		}
	}
	for _, v := range source.Defects {
		if (v.DetectionSprintID != nil && *v.DetectionSprintID == id) || (v.ResolutionSprintID != nil && *v.ResolutionSprintID == id) {
			result.Defects = append(result.Defects, v)
		}
	}
	return result, nil
}
