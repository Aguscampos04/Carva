package sprint

import "testing"

func TestNewSprint(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		goal      string
		wantGoal  string
		wantErr   bool
	}{
		{name: "trims goal", projectID: "project-1", goal: " MVP ", wantGoal: "MVP"},
		{name: "rejects invalid project", projectID: " ", goal: "MVP", wantErr: true},
		{name: "rejects blank goal", projectID: "project-1", goal: "  ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.projectID, tt.goal)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && (got.ProjectID != tt.projectID || got.Goal != tt.wantGoal) {
				t.Fatalf("New() = %+v, want project ID %q and goal %q", got, tt.projectID, tt.wantGoal)
			}
		})
	}
}
