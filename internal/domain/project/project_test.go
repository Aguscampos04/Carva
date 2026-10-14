package project

import (
	"errors"
	"testing"
)

func TestNewCreatesPlannedProjectWithStableID(t *testing.T) {
	got, err := New("project-123", "Carva MVP")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got.ID != "project-123" {
		t.Errorf("ID = %q, want %q", got.ID, "project-123")
	}
	if got.Name != "Carva MVP" {
		t.Errorf("Name = %q, want entered name unchanged", got.Name)
	}
	if got.Status != Planned {
		t.Errorf("Status = %q, want %q", got.Status, Planned)
	}
}

func TestNewRejectsNameThatIsEmptyAfterTrimming(t *testing.T) {
	for _, name := range []string{"", " ", "\t\r\n"} {
		t.Run(name, func(t *testing.T) {
			_, err := New("project-123", name)
			if !errors.Is(err, ErrInvalidName) {
				t.Errorf("New() error = %v, want ErrInvalidName", err)
			}
		})
	}
}

func TestNewRejectsEmptyID(t *testing.T) {
	_, err := New("", "Carva MVP")
	if !errors.Is(err, ErrInvalidID) {
		t.Errorf("New() error = %v, want ErrInvalidID", err)
	}
}
