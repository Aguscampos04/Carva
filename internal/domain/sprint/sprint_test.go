package sprint

import (
	"errors"
	"testing"
)

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

func TestSprintCanBeClosedOnlyOnce(t *testing.T) {
	current := Sprint{ID: 1, ProjectID: "project-1", Goal: "MVP"}
	closed, err := current.Close()
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !closed.Closed || current.Closed {
		t.Fatalf("Close() = %+v and original = %+v; want closed copy and unchanged original", closed, current)
	}
	if _, err := closed.Close(); !errors.Is(err, ErrAlreadyClosed) {
		t.Fatalf("Close() on closed Sprint error = %v, want ErrAlreadyClosed", err)
	}
}

func TestPlannedBacklogItemIsIncompleteUntilVerified(t *testing.T) {
	item := PlannedBacklogItem{AcceptanceCriteriaVerified: false}
	if !item.IsIncomplete() {
		t.Fatal("unverified planned story should be incomplete")
	}
	item.AcceptanceCriteriaVerified = true
	if item.IsIncomplete() {
		t.Fatal("verified planned story should be complete")
	}
}
