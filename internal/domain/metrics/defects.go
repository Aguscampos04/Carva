package metrics

import "github.com/Aguscampos04/Carva/internal/domain/defect"

type DefectSnapshot struct {
	Status             defect.Status
	DetectionSprintID  *int64
	ResolutionSprintID *int64
}

type DefectMetrics struct {
	DetectedDefects *int `json:"detected_defects"`
	ResolvedDefects *int `json:"resolved_defects"`
}

type SprintDefectMetrics struct {
	ID      int64         `json:"id"`
	Goal    string        `json:"goal"`
	Closed  bool          `json:"closed"`
	Metrics DefectMetrics `json:"metrics"`
}

type ProjectDefectMetrics struct {
	Metrics DefectMetrics         `json:"metrics"`
	Sprints []SprintDefectMetrics `json:"sprints"`
}

func CalculateProjectDefectMetrics(defects []DefectSnapshot) DefectMetrics {
	detected, resolved := 0, 0
	for _, item := range defects {
		detected++
		if isResolvedDefect(item.Status) {
			resolved++
		}
	}
	return DefectMetrics{DetectedDefects: intPointer(detected), ResolvedDefects: intPointer(resolved)}
}

func CalculateSprintDefectMetrics(sprintID int64, defects []DefectSnapshot) DefectMetrics {
	detected, resolved := 0, 0
	detectionAvailable, resolutionAvailable := true, true
	for _, item := range defects {
		if item.DetectionSprintID == nil {
			detectionAvailable = false
		} else if *item.DetectionSprintID == sprintID {
			detected++
		}
		if !isResolvedDefect(item.Status) {
			continue
		}
		if item.ResolutionSprintID == nil {
			resolutionAvailable = false
		} else if *item.ResolutionSprintID == sprintID {
			resolved++
		}
	}
	result := DefectMetrics{}
	if detectionAvailable {
		result.DetectedDefects = intPointer(detected)
	}
	if resolutionAvailable {
		result.ResolvedDefects = intPointer(resolved)
	}
	return result
}

func isResolvedDefect(status defect.Status) bool {
	return status == defect.StatusResolved || status == defect.StatusClosed
}
