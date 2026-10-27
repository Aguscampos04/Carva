// Package defect defines project-scoped defect records and their approved values.
package defect

import (
	"errors"
	"strings"
)

type Status string

type Severity string

const (
	StatusOpen       Status = "Abierto"
	StatusInProgress Status = "En progreso"
	StatusResolved   Status = "Resuelto"
	StatusClosed     Status = "Cerrado"
)

const (
	SeverityLow      Severity = "Baja"
	SeverityMedium   Severity = "Media"
	SeverityHigh     Severity = "Alta"
	SeverityCritical Severity = "Crítica"
)

var (
	ErrInvalidIdentity     = errors.New("defect identity, project, and backlog item are required")
	ErrInvalidDescription  = errors.New("defect description must not be blank")
	ErrInvalidSeverity     = errors.New("defect severity is invalid")
	ErrInvalidStatus       = errors.New("defect status is invalid")
	ErrInvalidSprint       = errors.New("defect Sprint IDs must be positive when provided")
	ErrProjectNotFound     = errors.New("defect project not found")
	ErrBacklogItemNotFound = errors.New("defect backlog item not found in project")
	ErrSprintNotFound      = errors.New("defect Sprint not found in project")
)

type Input struct {
	Description        string   `json:"description"`
	Severity           Severity `json:"severity"`
	Status             Status   `json:"status"`
	DetectionSprintID  *int64   `json:"detection_sprint_id"`
	ResolutionSprintID *int64   `json:"resolution_sprint_id"`
}

type Defect struct {
	ID                 string   `json:"id"`
	ProjectID          string   `json:"project_id"`
	BacklogItemID      string   `json:"backlog_item_id"`
	Description        string   `json:"description"`
	Severity           Severity `json:"severity"`
	Status             Status   `json:"status"`
	DetectionSprintID  *int64   `json:"detection_sprint_id"`
	ResolutionSprintID *int64   `json:"resolution_sprint_id"`
}

func New(id, projectID, backlogItemID string, input Input) (Defect, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(projectID) == "" || strings.TrimSpace(backlogItemID) == "" {
		return Defect{}, ErrInvalidIdentity
	}
	input.Description = strings.TrimSpace(input.Description)
	if input.Description == "" {
		return Defect{}, ErrInvalidDescription
	}
	switch input.Severity {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
	default:
		return Defect{}, ErrInvalidSeverity
	}
	switch input.Status {
	case StatusOpen, StatusInProgress, StatusResolved, StatusClosed:
	default:
		return Defect{}, ErrInvalidStatus
	}
	if input.DetectionSprintID != nil && *input.DetectionSprintID < 1 || input.ResolutionSprintID != nil && *input.ResolutionSprintID < 1 {
		return Defect{}, ErrInvalidSprint
	}
	return Defect{
		ID: id, ProjectID: projectID, BacklogItemID: backlogItemID,
		Description: input.Description, Severity: input.Severity, Status: input.Status,
		DetectionSprintID: cloneID(input.DetectionSprintID), ResolutionSprintID: cloneID(input.ResolutionSprintID),
	}, nil
}

func (input Input) Validate() error {
	_, err := New("defect", "project", "story", input)
	return err
}

func cloneID(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
