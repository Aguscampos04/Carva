package project

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidDate     = errors.New("invalid project date")
	ErrDateRange       = errors.New("start date is after finalization date")
	ErrStartAfterToday = errors.New("start date is after current finalization date")
	ErrFinalized       = errors.New("project is finalized")
	ErrEmptyPatch      = errors.New("at least one editable field is required")
)

// Date is a calendar day, not an instant. Its zero value means no recorded date.
type Date string

func ParseDate(value string) (Date, error) {
	day, err := time.Parse(time.DateOnly, value)
	if err != nil || day.Year() < 1 || day.Format(time.DateOnly) != value {
		return "", ErrInvalidDate
	}
	return Date(value), nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d == "" {
		return []byte("null"), nil
	}
	if _, err := ParseDate(string(d)); err != nil {
		return nil, err
	}
	return json.Marshal(string(d))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = ""
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return ErrInvalidDate
	}
	parsed, err := ParseDate(value)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// Patch distinguishes absence (keep) from an explicitly absent start (clear).
type Patch struct {
	Name         *string
	StartDate    Date
	HasStartDate bool
}

func (change Patch) Validate() error {
	if change.Name == nil && !change.HasStartDate {
		return ErrEmptyPatch
	}
	if change.Name != nil && strings.TrimSpace(*change.Name) == "" {
		return ErrInvalidName
	}
	if change.HasStartDate && change.StartDate != "" {
		_, err := ParseDate(string(change.StartDate))
		return err
	}
	return nil
}

// Edit returns a new value, leaving the receiver intact on all error paths.
func (p Project) Edit(change Patch) (Project, error) {
	if err := change.Validate(); err != nil {
		return Project{}, err
	}
	if change.Name != nil {
		p.Name = *change.Name
	}
	if change.HasStartDate {
		p.StartDate = change.StartDate
	}
	if p.StartDate != "" && p.EndDate != "" && p.StartDate > p.EndDate {
		return Project{}, ErrDateRange
	}
	return p, nil
}

// Finalize preserves the first end day, including unknown legacy end dates.
func (p Project) Finalize(today Date) (Project, error) {
	if p.Status == Finalized {
		return p, nil
	}
	if _, err := ParseDate(string(today)); err != nil {
		return Project{}, err
	}
	if p.StartDate != "" && p.StartDate > today {
		return Project{}, ErrStartAfterToday
	}
	p.Status = Finalized
	p.EndDate = today
	return p, nil
}

func (p Project) CanCreateSprint() error {
	if p.Status == Finalized {
		return ErrFinalized
	}
	return nil
}

// ValidateEffortDates protects active effort when project dates change (US-16 D-31).
var ErrEffortDates = errors.New("project dates conflict with active effort")

func (p Project) ValidateEffortDates(minimum, maximum Date) error {
	if minimum != "" && p.StartDate != "" && p.StartDate > minimum {
		return ErrEffortDates
	}
	if maximum != "" && p.Status == Finalized && p.EndDate != "" && maximum > p.EndDate {
		return ErrEffortDates
	}
	return nil
}
