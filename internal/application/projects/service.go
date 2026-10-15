// Package projects coordinates project creation and consultation use cases.
package projects

import (
	"context"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/google/uuid"
)

// Repository stores and retrieves the shared project entity.
type Repository interface {
	Save(context.Context, project.Project) error
	FindByID(context.Context, string) (project.Project, error)
}

// Service implements project creation and consultation use cases.
type Service struct {
	repository Repository
}

// NewService constructs a project service backed by repository.
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// Create validates and persists a project with a server-generated UUID.
func (s *Service) Create(ctx context.Context, name string) (project.Project, error) {
	created, err := project.New(uuid.NewString(), name)
	if err != nil {
		return project.Project{}, err
	}
	if err := s.repository.Save(ctx, created); err != nil {
		return project.Project{}, fmt.Errorf("save project: %w", err)
	}
	return created, nil
}

// Get retrieves a project by its stable identifier.
func (s *Service) Get(ctx context.Context, id string) (project.Project, error) {
	found, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return project.Project{}, fmt.Errorf("find project: %w", err)
	}
	return found, nil
}
