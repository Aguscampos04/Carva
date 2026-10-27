package metrics

import (
	"testing"

	"github.com/Aguscampos04/Carva/internal/domain/defect"
)

func TestCalculateProjectDefectMetrics(t *testing.T) {
	result := CalculateProjectDefectMetrics([]DefectSnapshot{
		{Status: defect.StatusOpen},
		{Status: defect.StatusResolved},
		{Status: defect.StatusClosed},
		{Status: defect.StatusInProgress},
	})
	if result.DetectedDefects == nil || *result.DetectedDefects != 4 || result.ResolvedDefects == nil || *result.ResolvedDefects != 2 {
		t.Fatalf("project defect metrics = %+v, want detected=4 and resolved=2", result)
	}
}

func TestCalculateSprintDefectMetricsSeparatesDetectionFromResolution(t *testing.T) {
	first, second := int64(1), int64(2)
	defects := []DefectSnapshot{
		{Status: defect.StatusResolved, DetectionSprintID: &first, ResolutionSprintID: &second},
		{Status: defect.StatusOpen, DetectionSprintID: &second},
		{Status: defect.StatusClosed, DetectionSprintID: &first, ResolutionSprintID: &first},
	}

	firstResult := CalculateSprintDefectMetrics(first, defects)
	if firstResult.DetectedDefects == nil || *firstResult.DetectedDefects != 2 || firstResult.ResolvedDefects == nil || *firstResult.ResolvedDefects != 1 {
		t.Fatalf("first Sprint metrics = %+v, want detected=2 and resolved=1", firstResult)
	}
	secondResult := CalculateSprintDefectMetrics(second, defects)
	if secondResult.DetectedDefects == nil || *secondResult.DetectedDefects != 1 || secondResult.ResolvedDefects == nil || *secondResult.ResolvedDefects != 1 {
		t.Fatalf("second Sprint metrics = %+v, want detected=1 and resolved=1", secondResult)
	}
}

func TestCalculateSprintDefectMetricsMarksMissingAttributionUnavailable(t *testing.T) {
	sprintID := int64(1)
	defects := []DefectSnapshot{
		{Status: defect.StatusOpen, DetectionSprintID: &sprintID},
		{Status: defect.StatusOpen},
		{Status: defect.StatusResolved, DetectionSprintID: &sprintID},
	}
	result := CalculateSprintDefectMetrics(sprintID, defects)
	if result.DetectedDefects != nil {
		t.Fatalf("detected defects = %v, want unavailable due to missing detection Sprint", *result.DetectedDefects)
	}
	if result.ResolvedDefects != nil {
		t.Fatal("resolved defects should be unavailable when a resolved defect lacks its resolution Sprint")
	}
}

func TestCalculateSprintDefectMetricsReturnsZeroWhenAttributionIsComplete(t *testing.T) {
	sprintID, otherID := int64(1), int64(2)
	result := CalculateSprintDefectMetrics(sprintID, []DefectSnapshot{
		{Status: defect.StatusOpen, DetectionSprintID: &otherID},
		{Status: defect.StatusResolved, DetectionSprintID: &otherID, ResolutionSprintID: &otherID},
	})
	if result.DetectedDefects == nil || *result.DetectedDefects != 0 || result.ResolvedDefects == nil || *result.ResolvedDefects != 0 {
		t.Fatalf("empty Sprint metrics = %+v, want available zero counts", result)
	}
}
