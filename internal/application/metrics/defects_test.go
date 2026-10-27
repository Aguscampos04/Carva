package metrics

import (
	"context"
	"testing"

	"github.com/Aguscampos04/Carva/internal/domain/defect"
	domain "github.com/Aguscampos04/Carva/internal/domain/metrics"
)

type defectMetricsRepository struct {
	data domain.ProjectDefectData
}

func (r defectMetricsRepository) ListProjectSprintSnapshots(context.Context, string) ([]domain.SprintSnapshot, error) {
	return nil, nil
}

func (r defectMetricsRepository) GetSprintSnapshot(context.Context, string, int64) (domain.SprintSnapshot, error) {
	return domain.SprintSnapshot{}, nil
}

func (r defectMetricsRepository) ProjectHoursData(context.Context, string) (domain.ProjectHoursData, error) {
	return domain.ProjectHoursData{}, nil
}

func (r defectMetricsRepository) ProjectDefectData(context.Context, string) (domain.ProjectDefectData, error) {
	return r.data, nil
}

func TestProjectDefectsBuildsProjectAndSprintMetrics(t *testing.T) {
	first, second := int64(1), int64(2)
	service := NewService(defectMetricsRepository{data: domain.ProjectDefectData{
		Sprints: []domain.DefectSprint{{ID: first, Goal: "First"}, {ID: second, Goal: "Second", Closed: true}},
		Defects: []domain.DefectSnapshot{
			{Status: defect.StatusResolved, DetectionSprintID: &first, ResolutionSprintID: &second},
			{Status: defect.StatusOpen, DetectionSprintID: &second},
		},
	}})

	result, err := service.ProjectDefects(context.Background(), "project")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Sprints) != 2 || *result.Metrics.DetectedDefects != 2 || *result.Metrics.ResolvedDefects != 1 {
		t.Fatalf("project defect metrics = %+v", result)
	}
	if *result.Sprints[0].Metrics.DetectedDefects != 1 || *result.Sprints[0].Metrics.ResolvedDefects != 0 ||
		*result.Sprints[1].Metrics.DetectedDefects != 1 || *result.Sprints[1].Metrics.ResolvedDefects != 1 {
		t.Fatalf("Sprint defect metrics = %+v", result.Sprints)
	}
}

func TestSprintDefectsRejectsInvalidOrMissingSprint(t *testing.T) {
	service := NewService(defectMetricsRepository{data: domain.ProjectDefectData{
		Sprints: []domain.DefectSprint{{ID: 1}},
	}})
	if _, err := service.SprintDefects(context.Background(), "project", 0); err == nil {
		t.Fatal("SprintDefects() accepted a non-positive Sprint ID")
	}
	if _, err := service.SprintDefects(context.Background(), "project", 2); err == nil {
		t.Fatal("SprintDefects() accepted a Sprint outside the project")
	}
}
