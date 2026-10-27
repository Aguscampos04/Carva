package defect

import (
	"errors"
	"testing"
)

func TestNewAcceptsEveryApprovedStatusAndSeverity(t *testing.T) {
	statuses := []Status{StatusOpen, StatusInProgress, StatusResolved, StatusClosed}
	severities := []Severity{SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical}
	for _, status := range statuses {
		for _, severity := range severities {
			t.Run(string(status)+"/"+string(severity), func(t *testing.T) {
				got, err := New("defect-id", "project-id", "story-id", Input{
					Description: "  No guarda cambios  ", Severity: severity, Status: status,
				})
				if err != nil {
					t.Fatal(err)
				}
				if got.Description != "No guarda cambios" || got.ProjectID != "project-id" || got.BacklogItemID != "story-id" || got.Severity != severity || got.Status != status {
					t.Fatalf("unexpected defect: %#v", got)
				}
				if got.DetectionSprintID != nil || got.ResolutionSprintID != nil {
					t.Fatalf("missing Sprints must remain absent: %#v", got)
				}
			})
		}
	}
}

func TestNewRejectsInvalidDefectFields(t *testing.T) {
	valid := Input{Description: "Bug", Severity: SeverityHigh, Status: StatusOpen}
	cases := []struct {
		name      string
		id        string
		projectID string
		storyID   string
		input     Input
		want      error
	}{
		{name: "missing identity", projectID: "p", storyID: "s", input: valid, want: ErrInvalidIdentity},
		{name: "blank description", id: "d", projectID: "p", storyID: "s", input: Input{Severity: SeverityHigh, Status: StatusOpen}, want: ErrInvalidDescription},
		{name: "unknown severity", id: "d", projectID: "p", storyID: "s", input: Input{Description: "Bug", Severity: "Urgente", Status: StatusOpen}, want: ErrInvalidSeverity},
		{name: "unknown status", id: "d", projectID: "p", storyID: "s", input: Input{Description: "Bug", Severity: SeverityHigh, Status: "Pendiente"}, want: ErrInvalidStatus},
		{name: "nonpositive detection Sprint", id: "d", projectID: "p", storyID: "s", input: Input{Description: "Bug", Severity: SeverityHigh, Status: StatusOpen, DetectionSprintID: int64Pointer(0)}, want: ErrInvalidSprint},
		{name: "nonpositive resolution Sprint", id: "d", projectID: "p", storyID: "s", input: Input{Description: "Bug", Severity: SeverityHigh, Status: StatusOpen, ResolutionSprintID: int64Pointer(-1)}, want: ErrInvalidSprint},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(tc.id, tc.projectID, tc.storyID, tc.input)
			if !errors.Is(err, tc.want) {
				t.Fatalf("New() error = %v, want %v", err, tc.want)
			}
		})
	}
}

func int64Pointer(value int64) *int64 { return &value }
