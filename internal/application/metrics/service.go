// Package metrics coordinates Story Points metric reads.
package metrics

import (
	"context"

	domain "github.com/Aguscampos04/Carva/internal/domain/metrics"
)

type Repository interface {
	ListProjectSprintSnapshots(context.Context, string) ([]domain.SprintSnapshot, error)
	GetSprintSnapshot(context.Context, string, int64) (domain.SprintSnapshot, error)
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
