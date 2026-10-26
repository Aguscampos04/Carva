package effort

import (
	"errors"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"testing"
)

// US-16 BDD 08/14/17/18/40/43: central rules independent of HTTP/SQLite.
func TestUS16EffortRules(t *testing.T) {
	p := project.Project{StartDate: "2026-10-01"}
	for _, tc := range []struct {
		date project.Date
		want error
	}{{"2026-10-01", nil}, {"2026-10-26", nil}, {"2026-09-30", ErrDateRange}, {"2026-10-27", ErrDateRange}} {
		if err := ValidateDate(p, "2026-10-26", tc.date); !errors.Is(err, tc.want) {
			t.Errorf("date %s: %v", tc.date, err)
		}
	}
	p.Status = project.Finalized
	if !errors.Is(ValidateDate(p, "2026-10-26", "2026-10-15"), ErrUnknownEnd) {
		t.Fatal("missing end accepted")
	}
	p.EndDate = "2026-10-20"
	if ValidateDate(p, "2026-10-26", "2026-10-20") != nil || !errors.Is(ValidateDate(p, "2026-10-26", "2026-10-21"), ErrDateRange) {
		t.Fatal("final end boundary")
	}
	in := Input{MemberID: "member", StoryID: "story", WorkDate: "2026-10-25", Activity: "review", Hours: 1}
	if in.Validate() != nil {
		t.Fatal("valid input rejected")
	}
	in.Activity = " \t "
	if !errors.Is(in.Validate(), ErrInvalid) {
		t.Fatal("blank accepted")
	}
}
