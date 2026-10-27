package defect

import "errors"

var (
	ErrEmptyPatch = errors.New("at least one editable defect field is required")
	ErrNotFound   = errors.New("defect not found in project and story")
)

// Patch distinguishes omitted Sprint associations from explicitly removed ones.
type Patch struct {
	Status              *Status
	HasDetectionSprint  bool
	DetectionSprintID   *int64
	HasResolutionSprint bool
	ResolutionSprintID  *int64
}

func (change Patch) Validate() error {
	if change.Status == nil && !change.HasDetectionSprint && !change.HasResolutionSprint {
		return ErrEmptyPatch
	}
	input := Input{Description: "validation", Severity: SeverityLow, Status: StatusOpen}
	if change.Status != nil {
		input.Status = *change.Status
	}
	if change.HasDetectionSprint {
		input.DetectionSprintID = change.DetectionSprintID
	}
	if change.HasResolutionSprint {
		input.ResolutionSprintID = change.ResolutionSprintID
	}
	return input.Validate()
}

// Update returns an independent copy, preserving identity and omitted fields.
func (value Defect) Update(change Patch) (Defect, error) {
	if err := change.Validate(); err != nil {
		return Defect{}, err
	}
	if change.Status != nil {
		value.Status = *change.Status
	}
	if change.HasDetectionSprint {
		value.DetectionSprintID = change.DetectionSprintID
	}
	if change.HasResolutionSprint {
		value.ResolutionSprintID = change.ResolutionSprintID
	}
	value.DetectionSprintID = cloneID(value.DetectionSprintID)
	value.ResolutionSprintID = cloneID(value.ResolutionSprintID)
	return value, nil
}
