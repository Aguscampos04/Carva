// Package backlog defines Product Backlog items and their business rules.
package backlog

import (
	"errors"
	"strings"
)

// Priority is a User Story's agreed Backlog priority.
type Priority string

const (
	PriorityHigh   Priority = "Alta"
	PriorityMedium Priority = "Media"
	PriorityLow    Priority = "Baja"
)

// Status is a User Story's lifecycle state in the Product Backlog.
type Status string

const (
	Pending    Status = "Pendiente"
	InProgress Status = "En progreso"
	Completed  Status = "Completada"
)

var (
	ErrInvalidID                     = errors.New("backlog item ID must not be empty")
	ErrInvalidProjectID              = errors.New("backlog project ID must not be empty")
	ErrInvalidTitle                  = errors.New("backlog item title must not be blank")
	ErrInvalidDescription            = errors.New("backlog item description must not be blank")
	ErrInvalidPriority               = errors.New("backlog item priority is invalid")
	ErrInvalidProvisionalStoryPoints = errors.New("provisional Story Points must be a non-negative integer")
	ErrInvalidAcceptanceCriteria     = errors.New("backlog item must have non-blank acceptance criteria")
	ErrEstimateFinalized             = errors.New("provisional Story Points cannot change after a definitive estimate exists")
	ErrNotFound                      = errors.New("backlog item not found")
	ErrProjectNotFound               = errors.New("backlog project not found")
)

// Input contains the fields users can edit through US-05.
type Input struct {
	Title                  string   `json:"title"`
	Description            string   `json:"description"`
	Priority               Priority `json:"priority"`
	ProvisionalStoryPoints *int     `json:"provisional_story_points"`
	AcceptanceCriteria     []string `json:"acceptance_criteria"`
}

// Item is a User Story stored in a project's Product Backlog.
type Item struct {
	ID                     string   `json:"id"`
	ProjectID              string   `json:"project_id"`
	Title                  string   `json:"title"`
	Description            string   `json:"description"`
	Priority               Priority `json:"priority"`
	Status                 Status   `json:"status"`
	ProvisionalStoryPoints *int     `json:"provisional_story_points"`
	StoryPoints            *int     `json:"story_points"`
	AcceptanceCriteria     []string `json:"acceptance_criteria"`
}

// New creates a pending Backlog item after validating its editable fields.
func New(id, projectID string, input Input) (Item, error) {
	if strings.TrimSpace(id) == "" {
		return Item{}, ErrInvalidID
	}
	if strings.TrimSpace(projectID) == "" {
		return Item{}, ErrInvalidProjectID
	}
	if err := validate(input); err != nil {
		return Item{}, err
	}

	return Item{
		ID:                     id,
		ProjectID:              projectID,
		Title:                  input.Title,
		Description:            input.Description,
		Priority:               input.Priority,
		Status:                 Pending,
		ProvisionalStoryPoints: cloneInt(input.ProvisionalStoryPoints),
		AcceptanceCriteria:     append([]string(nil), input.AcceptanceCriteria...),
	}, nil
}

// Update replaces editable fields while preserving identity, project, status,
// and any definitive estimate. US-05 cannot change provisional points after
// US-15 has recorded a definitive estimate.
func (item Item) Update(input Input) (Item, error) {
	if err := validate(input); err != nil {
		return Item{}, err
	}
	if item.StoryPoints != nil && !equalInt(item.ProvisionalStoryPoints, input.ProvisionalStoryPoints) {
		return Item{}, ErrEstimateFinalized
	}

	item.Title = input.Title
	item.Description = input.Description
	item.Priority = input.Priority
	item.ProvisionalStoryPoints = cloneInt(input.ProvisionalStoryPoints)
	item.AcceptanceCriteria = append([]string(nil), input.AcceptanceCriteria...)
	return item, nil
}

func validate(input Input) error {
	if strings.TrimSpace(input.Title) == "" {
		return ErrInvalidTitle
	}
	if strings.TrimSpace(input.Description) == "" {
		return ErrInvalidDescription
	}
	switch input.Priority {
	case PriorityHigh, PriorityMedium, PriorityLow:
	default:
		return ErrInvalidPriority
	}
	if input.ProvisionalStoryPoints != nil && *input.ProvisionalStoryPoints < 0 {
		return ErrInvalidProvisionalStoryPoints
	}
	if len(input.AcceptanceCriteria) == 0 {
		return ErrInvalidAcceptanceCriteria
	}
	for _, criterion := range input.AcceptanceCriteria {
		if strings.TrimSpace(criterion) == "" {
			return ErrInvalidAcceptanceCriteria
		}
	}
	return nil
}

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func equalInt(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
