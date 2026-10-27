package metrics

import (
	"math"
	"testing"
)

func points(value int) *int { return &value }

func TestCalculateSprintStoryPoints(t *testing.T) {
	tests := []struct {
		name   string
		closed bool
		items  []StorySnapshot
		want   StoryPointMetrics
	}{
		{
			name:   "complete points and completed percentage",
			closed: true,
			items:  []StorySnapshot{{StoryPoints: points(2), Completed: true}, {StoryPoints: points(3)}, {StoryPoints: points(5), Completed: true}},
			want:   StoryPointMetrics{PlannedStoryPoints: points(10), CompletedStoryPoints: points(7), TeamVelocity: floatPointer(7), CompletionPercentage: floatPointer(200.0 / 3.0)},
		},
		{
			name:  "open sprint has no velocity",
			items: []StorySnapshot{{StoryPoints: points(3)}},
			want:  StoryPointMetrics{PlannedStoryPoints: points(3), CompletedStoryPoints: points(0), CompletionPercentage: floatPointer(0)},
		},
		{
			name:   "missing story points only affect relevant sums",
			closed: true,
			items:  []StorySnapshot{{Completed: true}, {StoryPoints: points(4)}},
			want:   StoryPointMetrics{CompletedStoryPoints: nil, TeamVelocity: nil, CompletionPercentage: floatPointer(50)},
		},
		{
			name:   "empty closed sprint has zero sums but unavailable percentage",
			closed: true,
			want:   StoryPointMetrics{PlannedStoryPoints: points(0), CompletedStoryPoints: points(0), TeamVelocity: floatPointer(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateSprintStoryPoints(tt.closed, tt.items); !equalMetrics(got, tt.want) {
				t.Fatalf("CalculateSprintStoryPoints() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestCalculateProjectStoryPoints(t *testing.T) {
	tests := []struct {
		name    string
		sprints []SprintSnapshot
		want    StoryPointMetrics
	}{
		{
			name: "sums snapshots and weights completion while averaging closed sprint velocity",
			sprints: []SprintSnapshot{
				{Closed: true, Items: []StorySnapshot{{StoryPoints: points(3), Completed: true}, {StoryPoints: points(5)}}},
				{Closed: true, Items: []StorySnapshot{{StoryPoints: points(5), Completed: true}, {StoryPoints: points(2)}}},
				{Items: []StorySnapshot{{StoryPoints: points(8)}}},
			},
			want: StoryPointMetrics{PlannedStoryPoints: points(23), CompletedStoryPoints: points(8), TeamVelocity: floatPointer(4), CompletionPercentage: floatPointer(40)},
		},
		{
			name: "no sprint data makes project metrics unavailable",
			want: StoryPointMetrics{},
		},
		{
			name:    "closed sprint with incomplete completed score makes velocity unavailable",
			sprints: []SprintSnapshot{{Closed: true, Items: []StorySnapshot{{Completed: true}}}},
			want:    StoryPointMetrics{PlannedStoryPoints: nil, CompletedStoryPoints: nil, TeamVelocity: nil, CompletionPercentage: floatPointer(100)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateProjectStoryPoints(tt.sprints); !equalMetrics(got, tt.want) {
				t.Fatalf("CalculateProjectStoryPoints() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func equalMetrics(a, b StoryPointMetrics) bool {
	return equalInt(a.PlannedStoryPoints, b.PlannedStoryPoints) &&
		equalInt(a.CompletedStoryPoints, b.CompletedStoryPoints) &&
		equalFloat(a.TeamVelocity, b.TeamVelocity) &&
		equalFloat(a.CompletionPercentage, b.CompletionPercentage)
}

func equalInt(a, b *int) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func equalFloat(a, b *float64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return math.Abs(*a-*b) < 1e-9
}
