package sprint

import (
	"errors"
	"testing"
	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
)

func TestUS16SnapshotOwnsEstimateValue(t *testing.T) {
 item:=backlog.Item{ID:"a",EstimatedHours:hours.Ptr(800)}
 snap,err:=Snapshot(1,item);if err!=nil {t.Fatal(err)}
 *item.EstimatedHours=925
 if snap.EstimatedHours==nil || *snap.EstimatedHours!=800 {t.Fatal("snapshot shares mutable input estimate")}
}

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
