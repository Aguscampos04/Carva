package defect

import (
	"errors"
	"reflect"
	"testing"
)

func TestPatchPreservesIdentityAndOmittedFields(t *testing.T) {
	detection, resolution := int64(1), int64(2)
	original, err := New("d", "p", "a", Input{Description: "Bug", Severity: SeverityHigh, Status: StatusOpen, DetectionSprintID: &detection})
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []Status{StatusOpen, StatusInProgress, StatusResolved, StatusClosed} {
		got, err := original.Update(Patch{Status: &status, HasResolutionSprint: true, ResolutionSprintID: &resolution})
		if err != nil || got.Status != status || got.ResolutionSprintID == nil || *got.ResolutionSprintID != 2 || *got.DetectionSprintID != 1 {
			t.Fatalf("%+v %v", got, err)
		}
		if got.ID != original.ID || got.ProjectID != original.ProjectID || got.BacklogItemID != original.BacklogItemID || got.Description != original.Description || got.Severity != original.Severity {
			t.Fatal("noneditable fields changed")
		}
		*got.DetectionSprintID = 99
		if *original.DetectionSprintID != 1 || original.ResolutionSprintID != nil || original.Status != StatusOpen {
			t.Fatal("mutated source")
		}
	}
	got, err := original.Update(Patch{HasDetectionSprint: true})
	if err != nil || got.DetectionSprintID != nil || got.Status != original.Status {
		t.Fatalf("clear = %+v %v", got, err)
	}
}

func TestPatchRejectsInvalidValuesWithoutMutatingDefect(t *testing.T) {
	value := Defect{ID: "d", ProjectID: "p", BacklogItemID: "a", Description: "Bug", Severity: SeverityHigh, Status: StatusOpen}
	invalid := Status("Pendiente")
	zero, negative := int64(0), int64(-1)
	for _, tt := range []struct {
		change Patch
		want   error
	}{
		{Patch{}, ErrEmptyPatch},
		{Patch{Status: &invalid}, ErrInvalidStatus},
		{Patch{HasDetectionSprint: true, DetectionSprintID: &zero}, ErrInvalidSprint},
		{Patch{HasResolutionSprint: true, ResolutionSprintID: &negative}, ErrInvalidSprint},
	} {
		before := value
		if _, err := value.Update(tt.change); !errors.Is(err, tt.want) {
			t.Fatalf("error = %v, want %v", err, tt.want)
		}
		if !reflect.DeepEqual(value, before) {
			t.Fatal("mutated original on error")
		}
	}
}
