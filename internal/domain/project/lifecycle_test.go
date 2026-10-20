package project

import (
	"encoding/json"
	"errors"
	"testing"
)

// US-02 AC-02/03/06/08; BDD-03/05/06/09/10/14.
func TestCalendarDatesAndProjectEdits(t *testing.T) {
	for _, value := range []string{"2026-02-29", "2026-04-31", "2026-13-01", "2026-00-10", "2026-10-00", "20/10/2026", "2026-1-2", "2026-10-20T00:00:00Z", " 2026-10-01 ", "0000-01-01", "10000-01-01", ""} {
		if _, err := ParseDate(value); !errors.Is(err, ErrInvalidDate) {
			t.Errorf("ParseDate(%q) = %v", value, err)
		}
	}
	for _, value := range []string{"2024-02-29", "0001-01-01", "9999-12-31"} {
		if _, err := ParseDate(value); err != nil {
			t.Error(err)
		}
	}
	p, _ := New("id", "Carva Demo")
	name := "  Carva II  "
	p, err := p.Edit(Patch{Name: &name, HasStartDate: true, StartDate: Date("2026-10-01")})
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != name || p.StartDate != "2026-10-01" || p.Status != Planned {
		t.Fatalf("edit: %+v", p)
	}
	p, err = p.Finalize(Date("2026-10-20"))
	if err != nil {
		t.Fatal(err)
	}
	equal, err := p.Edit(Patch{HasStartDate: true, StartDate: Date("2026-10-20")})
	if err != nil || equal.EndDate != "2026-10-20" {
		t.Fatalf("equality: %+v %v", equal, err)
	}
	if _, err = p.Edit(Patch{Name: &name, HasStartDate: true, StartDate: Date("2026-10-21")}); !errors.Is(err, ErrDateRange) {
		t.Fatal(err)
	}
	cleared, err := p.Edit(Patch{HasStartDate: true})
	if err != nil || cleared.StartDate != "" || cleared.Status != Finalized || cleared.EndDate != p.EndDate {
		t.Fatalf("clear: %+v %v", cleared, err)
	}
	blank := " \t "
	if _, err = p.Edit(Patch{Name: &blank}); !errors.Is(err, ErrInvalidName) {
		t.Fatal(err)
	}
	if _, err = p.Edit(Patch{}); err == nil {
		t.Fatal("empty patch accepted")
	}
	if _, err = p.Edit(Patch{HasStartDate: true, StartDate: Date("impossible")}); !errors.Is(err, ErrInvalidDate) {
		t.Fatal(err)
	}
	raw, err := json.Marshal(cleared)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	_ = json.Unmarshal(raw, &decoded)
	if decoded["start_date"] != nil || decoded["end_date"] != "2026-10-20" {
		t.Fatalf("dates JSON: %s", raw)
	}
}

// US-02 AC-04/07/09/11; BDD-07/11/12/15/19/21.
func TestFinalizeIsManualAndRepeatable(t *testing.T) {
	for _, state := range []Status{Planned, InProgress} {
		p := Project{ID: "id", Name: "Demo", Status: state}
		got, err := p.Finalize(Date("2026-10-20"))
		if err != nil || got.Status != Finalized || got.EndDate != "2026-10-20" || got.StartDate != "" {
			t.Fatalf("%+v %v", got, err)
		}
		repeated, err := got.Finalize(Date("2026-10-21"))
		if err != nil || repeated != got {
			t.Fatalf("repeat: %+v %v", repeated, err)
		}
		if !errors.Is(got.CanCreateSprint(), ErrFinalized) || p.CanCreateSprint() != nil {
			t.Fatal("Sprint eligibility")
		}
	}
	p := Project{Status: Planned, StartDate: Date("2026-10-21")}
	if _, err := p.Finalize(Date("2026-10-20")); !errors.Is(err, ErrStartAfterToday) {
		t.Fatalf("future: %v", err)
	}
	if p.Status != Planned || p.EndDate != "" {
		t.Fatal("mutated receiver")
	}
	if _, err := p.Finalize(Date("")); !errors.Is(err, ErrInvalidDate) {
		t.Fatal(err)
	}
	legacy := Project{Status: Finalized}
	got, err := legacy.Finalize(Date("2026-10-20"))
	if err != nil || got.EndDate != "" {
		t.Fatal("invented legacy date")
	}
}
