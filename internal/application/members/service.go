package members

import (
	"context"
	"fmt"
	"strings"

	"github.com/Aguscampos04/Carva/internal/domain/member"
	"github.com/google/uuid"
)

// Repository stores registrations and checks their owning project.
type Repository interface {
	ProjectExists(context.Context, string) (bool, error)
	Create(context.Context, member.Member) error
	ListByProject(context.Context, string) ([]member.Member, error)
}

// Service coordinates member registration and project-scoped listing.
type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Create(ctx context.Context, projectID, name string) (member.Member, error) {
	if strings.TrimSpace(name) == "" {
		return member.Member{}, member.ErrInvalidName
	}
	exists, err := s.repository.ProjectExists(ctx, projectID)
	if err != nil {
		return member.Member{}, fmt.Errorf("check member project: %w", err)
	}
	if !exists {
		return member.Member{}, member.ErrProjectNotFound
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return member.Member{}, fmt.Errorf("generate member UUID: %w", err)
	}
	created, err := member.New(id.String(), projectID, name)
	if err != nil {
		return member.Member{}, err
	}
	if err := s.repository.Create(ctx, created); err != nil {
		return member.Member{}, fmt.Errorf("register member: %w", err)
	}
	return created, nil
}

func (s *Service) List(ctx context.Context, projectID string) ([]member.Member, error) {
	exists, err := s.repository.ProjectExists(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("check member project: %w", err)
	}
	if !exists {
		return nil, member.ErrProjectNotFound
	}
	items, err := s.repository.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	if items == nil {
		items = []member.Member{}
	}
	return items, nil
}
