// Package metrics contains project and Sprint metric calculations.
package metrics

type StorySnapshot struct {
	StoryPoints *int
	Completed   bool
}

type SprintSnapshot struct {
	ID     int64
	Goal   string
	Closed bool
	Items  []StorySnapshot
}

type StoryPointMetrics struct {
	PlannedStoryPoints   *int     `json:"planned_story_points"`
	CompletedStoryPoints *int     `json:"completed_story_points"`
	TeamVelocity         *float64 `json:"team_velocity"`
	CompletionPercentage *float64 `json:"completion_percentage"`
}

type SprintMetrics struct {
	ID      int64             `json:"id"`
	Goal    string            `json:"goal"`
	Closed  bool              `json:"closed"`
	Metrics StoryPointMetrics `json:"metrics"`
}

type ProjectStoryPointMetrics struct {
	Metrics StoryPointMetrics `json:"metrics"`
	Sprints []SprintMetrics   `json:"sprints"`
}

func CalculateSprintStoryPoints(closed bool, items []StorySnapshot) StoryPointMetrics {
	planned, completed := 0, 0
	plannedMissing, completedMissing := false, false
	completedCount := 0

	for _, item := range items {
		if item.StoryPoints == nil {
			plannedMissing = true
			if item.Completed {
				completedMissing = true
			}
		} else {
			planned += *item.StoryPoints
			if item.Completed {
				completed += *item.StoryPoints
			}
		}
		if item.Completed {
			completedCount++
		}
	}

	result := StoryPointMetrics{}
	if !plannedMissing {
		result.PlannedStoryPoints = intPointer(planned)
	}
	if !completedMissing {
		result.CompletedStoryPoints = intPointer(completed)
	}
	if closed && result.CompletedStoryPoints != nil {
		result.TeamVelocity = floatPointer(float64(completed))
	}
	if len(items) > 0 {
		result.CompletionPercentage = floatPointer(float64(completedCount) / float64(len(items)) * 100)
	}
	return result
}

func CalculateProjectStoryPoints(sprints []SprintSnapshot) StoryPointMetrics {
	if len(sprints) == 0 {
		return StoryPointMetrics{}
	}

	plannedTotal, completedTotal := 0, 0
	plannedAvailable, completedAvailable := true, true
	completedPointsByClosedSprint := 0.0
	closedSprints, velocityAvailable := 0, true
	totalStories, completedStories := 0, 0

	for _, current := range sprints {
		value := CalculateSprintStoryPoints(current.Closed, current.Items)
		if value.PlannedStoryPoints == nil {
			plannedAvailable = false
		} else {
			plannedTotal += *value.PlannedStoryPoints
		}
		if value.CompletedStoryPoints == nil {
			completedAvailable = false
		} else {
			completedTotal += *value.CompletedStoryPoints
		}
		if current.Closed {
			closedSprints++
			if value.TeamVelocity == nil {
				velocityAvailable = false
			} else {
				completedPointsByClosedSprint += *value.TeamVelocity
			}
		}
		for _, item := range current.Items {
			totalStories++
			if item.Completed {
				completedStories++
			}
		}
	}

	result := StoryPointMetrics{}
	if plannedAvailable {
		result.PlannedStoryPoints = intPointer(plannedTotal)
	}
	if completedAvailable {
		result.CompletedStoryPoints = intPointer(completedTotal)
	}
	if closedSprints > 0 && velocityAvailable {
		result.TeamVelocity = floatPointer(completedPointsByClosedSprint / float64(closedSprints))
	}
	if totalStories > 0 {
		result.CompletionPercentage = floatPointer(float64(completedStories) / float64(totalStories) * 100)
	}
	return result
}

func intPointer(value int) *int           { return &value }
func floatPointer(value float64) *float64 { return &value }
