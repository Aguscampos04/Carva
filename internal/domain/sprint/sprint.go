// Package sprint contains Sprint entities and business rules.
package sprint

import (
	"errors"
	"strings"
)

var (
	ErrInvalidProjectID = errors.New("project ID must not be blank")
	ErrInvalidGoal      = errors.New("sprint goal must not be blank")
	ErrNotFound         = errors.New("sprint not found")
)

type Sprint struct {
	ID        int64  `json:"id"`
	ProjectID string `json:"projectID"`
	Goal      string `json:"goal"`
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
