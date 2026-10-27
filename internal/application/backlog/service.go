// Package backlog coordinates Product Backlog use cases.
package backlog

import (
	"context"
	"fmt"

	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/google/uuid"
)

// Repository stores Backlog items and checks their owning project.
type Repository interface {
	ProjectExists(context.Context, string) (bool, error)
	Create(context.Context, domain.Item) error
	FindByID(context.Context, string, string) (domain.Item, error)
	Mutate(context.Context, string, string, func(domain.Item) (domain.Item, error)) (domain.Item, error)
	ListByProject(context.Context, string) ([]domain.Item, error)
}

// Service implements the Product Backlog create, update, and list use cases.
type Service struct {
	repository Repository
}

// NewService constructs a Backlog service backed by repository.
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// Create validates and persists an item with a server-generated UUID.
func (s *Service) Create(ctx context.Context, projectID string, input domain.Input) (domain.Item, error) {
	exists, err := s.repository.ProjectExists(ctx, projectID)
	if err != nil {
		return domain.Item{}, fmt.Errorf("check backlog project: %w", err)
	}
	if !exists {
		return domain.Item{}, domain.ErrProjectNotFound
	}

	created, err := domain.New(uuid.NewString(), projectID, input)
	if err != nil {
		return domain.Item{}, err
	}
	if err := s.repository.Create(ctx, created); err != nil {
		return domain.Item{}, fmt.Errorf("create backlog item: %w", err)
	}
	return created, nil
}

// Update replaces editable fields while preserving controlled values.
func (s *Service) Update(ctx context.Context, projectID, itemID string, input domain.Input) (domain.Item, error) {
	return s.repository.Mutate(ctx, projectID, itemID, func(current domain.Item) (domain.Item, error) { return current.Update(input) })
}

// List returns all Backlog items belonging to an existing project.
func (s *Service) List(ctx context.Context, projectID string) ([]domain.Item, error) {
	exists, err := s.repository.ProjectExists(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("check backlog project: %w", err)
	}
	if !exists {
		return nil, domain.ErrProjectNotFound
	}

	items, err := s.repository.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list backlog items: %w", err)
	}
	if items == nil {
		return []domain.Item{}, nil
	}
	return items, nil
}

func (s *Service) Get(ctx context.Context, pid, id string) (domain.Item, error) {
	return s.repository.FindByID(ctx, pid, id)
}
