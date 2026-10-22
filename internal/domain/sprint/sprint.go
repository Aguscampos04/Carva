// Package sprint contains Sprint entities and business rules.
package sprint

import (
	"errors"
	"strings"

	"github.com/Aguscampos04/Carva/internal/domain/backlog"
)

var (
	ErrInvalidProjectID           = errors.New("project ID must not be blank")
	ErrInvalidGoal                = errors.New("sprint goal must not be blank")
	ErrInvalidBacklogItemID       = errors.New("backlog item ID must not be blank")
	ErrBacklogItemAlreadyAssigned = errors.New("backlog item is already assigned to this sprint")
	ErrBacklogItemNotAssigned     = errors.New("backlog item is not assigned to this sprint")
	ErrNotFound                   = errors.New("sprint not found")
	ErrAlreadyClosed              = errors.New("sprint is already closed")
	ErrClosed                     = errors.New("sprint is closed")
	ErrNextSprintRequired         = errors.New("an open subsequent sprint is required for incomplete stories")
)

type Sprint struct {
	ID        int64  `json:"id"`
	ProjectID string `json:"projectID"`
	Goal      string `json:"goal"`
	Closed    bool   `json:"closed"`
}

// Close returns a closed copy and rejects attempts to close the same Sprint twice.
func (value Sprint) Close() (Sprint, error) {
	if value.Closed {
		return Sprint{}, ErrAlreadyClosed
	}
	value.Closed = true
	return value, nil
}

// PlannedBacklogItem is the immutable planning-time copy of a Backlog item.
type PlannedBacklogItem struct {
	SprintID                   int64    `json:"sprint_id"`
	BacklogItemID              string   `json:"backlog_item_id"`
	Title                      string   `json:"title"`
	Description                string   `json:"description"`
	Priority                   string   `json:"priority"`
	Status                     string   `json:"status"`
	PlannedStatus              string   `json:"planned_status"`
	AcceptanceCriteriaVerified bool     `json:"acceptance_criteria_verified"`
	ProvisionalStoryPoints     *int     `json:"provisional_story_points"`
	StoryPoints                *int     `json:"story_points"`
	AcceptanceCriteria         []string `json:"acceptance_criteria"`
}

// Snapshot captures a Backlog item's planning fields for a Sprint.
func Snapshot(sprintID int64, item backlog.Item) (PlannedBacklogItem, error) {
	if sprintID < 1 {
		return PlannedBacklogItem{}, ErrNotFound
	}
	if strings.TrimSpace(item.ID) == "" {
		return PlannedBacklogItem{}, ErrInvalidBacklogItemID
	}
	return PlannedBacklogItem{
		SprintID:               sprintID,
		BacklogItemID:          item.ID,
		Title:                  item.Title,
		Description:            item.Description,
		Priority:               string(item.Priority),
		Status:                 string(item.Status),
		PlannedStatus:          string(item.Status),
		ProvisionalStoryPoints: cloneInt(item.ProvisionalStoryPoints),
		StoryPoints:            cloneInt(item.StoryPoints),
		AcceptanceCriteria:     append([]string(nil), item.AcceptanceCriteria...),
	}, nil
}

func New(projectID string, goal string) (Sprint, error) {
	if strings.TrimSpace(projectID) == "" {
		return Sprint{}, ErrInvalidProjectID
	}
	goal = strings.TrimSpace(goal)
	if goal == "" {
		return Sprint{}, ErrInvalidGoal
	}
	return Sprint{ProjectID: projectID, Goal: goal}, nil
}

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
