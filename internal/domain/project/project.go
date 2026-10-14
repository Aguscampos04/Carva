// Package project defines the project entity and its creation rules.
package project

import (
	"errors"
	"strings"
)

// Status is a project's lifecycle state.
type Status string

const (
	Planned    Status = "Planificado"
	InProgress Status = "En progreso"
	Finalized  Status = "Finalizado"
)

var (
	ErrInvalidID   = errors.New("project ID must not be empty")
	ErrInvalidName = errors.New("project name must not be blank")
	ErrNotFound    = errors.New("project not found")
)

// Project is the shared project model used by project management features.
type Project struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

// New creates a project with the initial state required by US-01.
func New(id, name string) (Project, error) {
	if strings.TrimSpace(id) == "" {
		return Project{}, ErrInvalidID
	}
	if strings.TrimSpace(name) == "" {
		return Project{}, ErrInvalidName
	}

	return Project{ID: id, Name: name, Status: Planned}, nil
}
