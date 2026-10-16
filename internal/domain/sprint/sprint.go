// Package sprint contains Sprint entities and business rules.
package sprint

import (
	"errors"
	"strings"
)

var (
	ErrInvalidProjectID = errors.New("project ID must be positive")
	ErrInvalidGoal      = errors.New("sprint goal must not be blank")
	ErrNotFound         = errors.New("sprint not found")
)

type Sprint struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"projectID"`
	Goal      string `json:"goal"`
}

func New(projectID int64, goal string) (Sprint, error) {
	if projectID < 1 {
		return Sprint{}, ErrInvalidProjectID
	}
	goal = strings.TrimSpace(goal)
	if goal == "" {
		return Sprint{}, ErrInvalidGoal
	}
	return Sprint{ProjectID: projectID, Goal: goal}, nil
}
