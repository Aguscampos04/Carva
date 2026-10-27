// Package defects coordinates the registration and consultation of defects.
package defects

import (
	"context"
	"fmt"

	domain "github.com/Aguscampos04/Carva/internal/domain/defect"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, domain.Defect) error
	ListByStory(context.Context, string, string) ([]domain.Defect, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Create(ctx context.Context, projectID, backlogItemID string, input domain.Input) (domain.Defect, error) {
	if err := input.Validate(); err != nil {
		return domain.Defect{}, err
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return domain.Defect{}, fmt.Errorf("generate defect UUID: %w", err)
	}
	created, err := domain.New(id.String(), projectID, backlogItemID, input)
	if err != nil {
		return domain.Defect{}, err
	}
	if err := s.repository.Create(ctx, created); err != nil {
		return domain.Defect{}, fmt.Errorf("register defect: %w", err)
	}
	return created, nil
}

func (s *Service) ListByStory(ctx context.Context, projectID, backlogItemID string) ([]domain.Defect, error) {
	items, err := s.repository.ListByStory(ctx, projectID, backlogItemID)
	if err != nil {
		return nil, fmt.Errorf("list defects by story: %w", err)
	}
	if items == nil {
		return []domain.Defect{}, nil
	}
	return items, nil
}
