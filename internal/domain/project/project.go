// Package project contains project entities and business rules.
package project

import (
	"errors"
	"strings"
)

type Status string

const (
	Planned    Status = "Planificado"
	InProgress Status = "En progreso"
	Finalized  Status = "Finalizado"
)

var (
	ErrInvalidName = errors.New("project name must not be blank")
	ErrInvalidID   = errors.New("project ID must be positive")
	ErrNotFound    = errors.New("project not found")
)

type Project struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

func New(name string) (Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Project{}, ErrInvalidName
	}
	return Project{Name: name, Status: Planned}, nil
}

// StartFirstSprint applies the approved transition when a planned project starts work.
func (p *Project) StartFirstSprint() {
	if p.Status == Planned {
		p.Status = InProgress
	}
}
