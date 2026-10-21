package backlog

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewCreatesPendingBacklogItemWithIndependentStoryPointValues(t *testing.T) {
	provisional := 5
	got, err := New("item-1", "project-1", Input{
		Title:                  "  Login  ",
		Description:            "  User can sign in  ",
		Priority:               PriorityHigh,
		ProvisionalStoryPoints: &provisional,
		AcceptanceCriteria:     []string{"Valid credentials grant access"},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got.ID != "item-1" || got.ProjectID != "project-1" || got.Title != "  Login  " || got.Description != "  User can sign in  " {
		t.Errorf("New() did not preserve identity and text: %+v", got)
	}
	if got.Status != Pending {
		t.Errorf("New() status = %q, want %q", got.Status, Pending)
	}
	if got.ProvisionalStoryPoints == nil || *got.ProvisionalStoryPoints != provisional {
		t.Errorf("New() provisional points = %v, want %d", got.ProvisionalStoryPoints, provisional)
	}
	if got.StoryPoints != nil {
		t.Errorf("New() definitive points = %v, want nil", got.StoryPoints)
	}
	if !reflect.DeepEqual(got.AcceptanceCriteria, []string{"Valid credentials grant access"}) {
		t.Errorf("New() acceptance criteria = %#v", got.AcceptanceCriteria)
	}
}

func TestNewAllowsMissingProvisionalStoryPoints(t *testing.T) {
	got, err := New("item-1", "project-1", validInput())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got.ProvisionalStoryPoints != nil || got.StoryPoints != nil {
		t.Errorf("New() points = (%v, %v), want both nil", got.ProvisionalStoryPoints, got.StoryPoints)
	}
}

func TestNewRejectsInvalidRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		id    string
		pid   string
		input Input
		want  error
	}{
		{name: "missing id", id: " ", pid: "project-1", input: validInput(), want: ErrInvalidID},
		{name: "missing project", id: "item-1", pid: " ", input: validInput(), want: ErrInvalidProjectID},
		{name: "blank title", id: "item-1", pid: "project-1", input: func() Input { in := validInput(); in.Title = " \t "; return in }(), want: ErrInvalidTitle},
		{name: "blank description", id: "item-1", pid: "project-1", input: func() Input { in := validInput(); in.Description = " \n "; return in }(), want: ErrInvalidDescription},
		{name: "unknown priority", id: "item-1", pid: "project-1", input: func() Input { in := validInput(); in.Priority = "Urgente"; return in }(), want: ErrInvalidPriority},
		{name: "missing criteria", id: "item-1", pid: "project-1", input: func() Input { in := validInput(); in.AcceptanceCriteria = nil; return in }(), want: ErrInvalidAcceptanceCriteria},
		{name: "blank criterion", id: "item-1", pid: "project-1", input: func() Input { in := validInput(); in.AcceptanceCriteria = []string{"  "}; return in }(), want: ErrInvalidAcceptanceCriteria},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.id, tt.pid, tt.input)
			if !errors.Is(err, tt.want) {
				t.Errorf("New() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestNewRejectsNegativeProvisionalStoryPoints(t *testing.T) {
	negative := -1
	in := validInput()
	in.ProvisionalStoryPoints = &negative
	_, err := New("item-1", "project-1", in)
	if !errors.Is(err, ErrInvalidProvisionalStoryPoints) {
		t.Errorf("New() error = %v, want ErrInvalidProvisionalStoryPoints", err)
	}
}

func TestUpdateChangesEditableFieldsAndPreservesControlledFields(t *testing.T) {
	item, err := New("item-1", "project-1", validInput())
	if err != nil {
		t.Fatal(err)
	}
	item.Status = InProgress
	definitive := 8
	item.StoryPoints = &definitive
	updatedInput := Input{
		Title: "Updated title", Description: "Updated description", Priority: PriorityLow,
		AcceptanceCriteria: []string{"Updated criterion"},
	}
	got, err := item.Update(updatedInput)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if got.ID != item.ID || got.ProjectID != item.ProjectID || got.Status != item.Status {
		t.Errorf("Update() changed controlled fields: %+v", got)
	}
	if got.StoryPoints == nil || *got.StoryPoints != definitive {
		t.Errorf("Update() definitive points = %v, want %d", got.StoryPoints, definitive)
	}
	if got.Title != updatedInput.Title || got.Description != updatedInput.Description || got.Priority != updatedInput.Priority {
		t.Errorf("Update() editable fields = %+v", got)
	}
}

func TestUpdateRejectsProvisionalEditAfterDefinitiveEstimateExists(t *testing.T) {
	item, err := New("item-1", "project-1", validInput())
	if err != nil {
		t.Fatal(err)
	}
	definitive := 8
	item.StoryPoints = &definitive
	provisional := 3
	in := validInput()
	in.ProvisionalStoryPoints = &provisional
	if _, err := item.Update(in); !errors.Is(err, ErrEstimateFinalized) {
		t.Errorf("Update() error = %v, want ErrEstimateFinalized", err)
	}
}

func TestCompleteRequiresExplicitAcceptanceCriteriaVerification(t *testing.T) {
	item, err := New("item-1", "project-1", validInput())
	if err != nil {
		t.Fatal(err)
	}
	item.Status = InProgress

	if _, err := item.Complete(false); !errors.Is(err, ErrAcceptanceCriteriaUnverified) {
		t.Fatalf("Complete(false) error = %v, want ErrAcceptanceCriteriaUnverified", err)
	}
	if item.Status != InProgress {
		t.Fatalf("item status after rejected completion = %q, want %q", item.Status, InProgress)
	}

	completed, err := item.Complete(true)
	if err != nil {
		t.Fatalf("Complete(true) error = %v", err)
	}
	if completed.Status != Completed || item.Status != InProgress {
		t.Fatalf("Complete(true) status = %q and original = %q; want completed copy and unchanged original", completed.Status, item.Status)
	}
}

func validInput() Input {
	return Input{
		Title: "Login", Description: "User can sign in", Priority: PriorityMedium,
		AcceptanceCriteria: []string{"Valid credentials grant access"},
	}
}
