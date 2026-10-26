// Package effort defines US-16 effort rules without transport or storage dependencies.
package effort

import (
	"errors"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"strings"
)

var (
	ErrInvalid    = errors.New("member, story, date, activity and positive hours are required")
	ErrNotFound   = errors.New("effort resource not found in project")
	ErrSprint     = errors.New("select a Sprint belonging to the destination story or explicitly remove the association")
	ErrDateRange  = errors.New("work date is outside the project dates or after today")
	ErrUnknownEnd = errors.New("cannot register or edit effort: finalized project has unknown end date")
	ErrDailyLimit = errors.New("member cannot exceed 24 hours per project and work date")
)

type Input struct {
	MemberID  string       `json:"member_id"`
	StoryID   string       `json:"backlog_item_id"`
	WorkDate  project.Date `json:"work_date"`
	Activity  string       `json:"activity"`
	Hours     hours.Value  `json:"hours"`
	SprintID  *int64       `json:"sprint_id"`
	HasSprint bool         `json:"-"`
}
type Entry struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Input
	DeletedAt *string `json:"-"`
}

func (in Input) Validate() error {
	if strings.TrimSpace(in.MemberID) == "" || strings.TrimSpace(in.StoryID) == "" || strings.TrimSpace(in.Activity) == "" || in.Hours <= 0 {
		return ErrInvalid
	}
	if _, err := project.ParseDate(string(in.WorkDate)); err != nil {
		return ErrInvalid
	}
	if in.SprintID != nil && *in.SprintID < 1 {
		return ErrInvalid
	}
	return nil
}
func ValidateDate(p project.Project, today, day project.Date) error {
	if err := ValidateWritable(p); err != nil {
		return err
	}
	if day > today || (p.StartDate != "" && day < p.StartDate) || (p.Status == project.Finalized && day > p.EndDate) {
		return ErrDateRange
	}
	return nil
}

// ValidateWritable is also checked before relationship lookup, to give D-30
// the same precedence in all effort writes. Consultation/deletion do not call it.
func ValidateWritable(p project.Project) error {
	if p.Status == project.Finalized && p.EndDate == "" {
		return ErrUnknownEnd
	}
	return nil
}
func ValidateDaily(other, added hours.Value) error {
	if added <= 0 || added > 2400 || other > 2400-added {
		return ErrDailyLimit
	}
	return nil
}
