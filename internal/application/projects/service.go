package projects

import (
	"context"

	"github.com/Aguscampos04/Carva/internal/domain/project"
)

type Repository interface {
	Create(context.Context, project.Project) (project.Project, error)
	Get(context.Context, int64) (project.Project, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) Service { return Service{repository: repository} }

func (s Service) Create(ctx context.Context, name string) (project.Project, error) {
	value, err := project.New(name)
	if err != nil {
		return project.Project{}, err
	}
	return s.repository.Create(ctx, value)
}

func (s Service) Get(ctx context.Context, id int64) (project.Project, error) {
	if id < 1 {
		return project.Project{}, project.ErrInvalidID
	}
	return s.repository.Get(ctx, id)
}
