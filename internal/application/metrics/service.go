// Package metrics coordinates Story Points metric reads.
package metrics

import (
	"context"

	domain "github.com/Aguscampos04/Carva/internal/domain/metrics"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type Repository interface {
	ListProjectSprintSnapshots(context.Context, string) ([]domain.SprintSnapshot, error)
	GetSprintSnapshot(context.Context, string, int64) (domain.SprintSnapshot, error)
	ProjectHoursData(context.Context, string) (domain.ProjectHoursData, error)
	ProjectDefectData(context.Context, string) (domain.ProjectDefectData, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) ProjectStoryPoints(ctx context.Context, projectID string) (domain.ProjectStoryPointMetrics, error) {
	snapshots, err := s.repository.ListProjectSprintSnapshots(ctx, projectID)
	if err != nil {
		return domain.ProjectStoryPointMetrics{}, err
	}

	result := domain.ProjectStoryPointMetrics{
		Metrics: domain.CalculateProjectStoryPoints(snapshots),
		Sprints: make([]domain.SprintMetrics, 0, len(snapshots)),
	}
	for _, snapshot := range snapshots {
		result.Sprints = append(result.Sprints, domain.SprintMetrics{
			ID: snapshot.ID, Goal: snapshot.Goal, Closed: snapshot.Closed,
			Metrics: domain.CalculateSprintStoryPoints(snapshot.Closed, snapshot.Items),
		})
	}
	return result, nil
}

func (s *Service) SprintStoryPoints(ctx context.Context, projectID string, sprintID int64) (domain.StoryPointMetrics, error) {
	snapshot, err := s.repository.GetSprintSnapshot(ctx, projectID, sprintID)
	if err != nil {
		return domain.StoryPointMetrics{}, err
	}
	return domain.CalculateSprintStoryPoints(snapshot.Closed, snapshot.Items), nil
}

func (s *Service) ProjectHours(ctx context.Context, projectID string) (domain.ProjectHoursMetrics, error) {
	data, err := s.repository.ProjectHoursData(ctx, projectID)
	if err != nil {
		return domain.ProjectHoursMetrics{}, err
	}
	result := domain.ProjectHoursMetrics{
		Sprints: make([]domain.SprintHoursMetrics, 0, len(data.Sprints)),
	}
	result.Metrics, err = domain.CalculateHours(data.Estimates, data.ActualHours)
	if err != nil {
		return domain.ProjectHoursMetrics{}, err
	}
	for _, sprintData := range data.Sprints {
		metrics, err := domain.CalculateHours(sprintData.Estimates, sprintData.ActualHours)
		if err != nil {
			return domain.ProjectHoursMetrics{}, err
		}
		result.Sprints = append(result.Sprints, domain.SprintHoursMetrics{
			ID: sprintData.ID, Goal: sprintData.Goal, Closed: sprintData.Closed, Metrics: metrics,
		})
	}
	return result, nil
}

func (s *Service) SprintHours(ctx context.Context, projectID string, sprintID int64) (domain.HoursMetrics, error) {
	if sprintID < 1 {
		return domain.HoursMetrics{}, sprint.ErrNotFound
	}
	data, err := s.repository.ProjectHoursData(ctx, projectID)
	if err != nil {
		return domain.HoursMetrics{}, err
	}
	for _, sprintData := range data.Sprints {
		if sprintData.ID == sprintID {
			return domain.CalculateHours(sprintData.Estimates, sprintData.ActualHours)
		}
	}
	return domain.HoursMetrics{}, sprint.ErrNotFound
}

func (s *Service) ProjectDefects(ctx context.Context, projectID string) (domain.ProjectDefectMetrics, error) {
	data, err := s.repository.ProjectDefectData(ctx, projectID)
	if err != nil {
		return domain.ProjectDefectMetrics{}, err
	}
	result := domain.ProjectDefectMetrics{
		Metrics: domain.CalculateProjectDefectMetrics(data.Defects),
		Sprints: make([]domain.SprintDefectMetrics, 0, len(data.Sprints)),
	}
	for _, sprintData := range data.Sprints {
		result.Sprints = append(result.Sprints, domain.SprintDefectMetrics{
			ID: sprintData.ID, Goal: sprintData.Goal, Closed: sprintData.Closed,
			Metrics: domain.CalculateSprintDefectMetrics(sprintData.ID, data.Defects),
		})
	}
	return result, nil
}

func (s *Service) SprintDefects(ctx context.Context, projectID string, sprintID int64) (domain.DefectMetrics, error) {
	if sprintID < 1 {
		return domain.DefectMetrics{}, sprint.ErrNotFound
	}
	data, err := s.repository.ProjectDefectData(ctx, projectID)
	if err != nil {
		return domain.DefectMetrics{}, err
	}
	for _, value := range data.Sprints {
		if value.ID == sprintID {
			return domain.CalculateSprintDefectMetrics(sprintID, data.Defects), nil
		}
	}
	return domain.DefectMetrics{}, sprint.ErrNotFound
}
