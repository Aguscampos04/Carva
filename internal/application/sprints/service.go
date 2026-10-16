package sprints

import (
	"context"

	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type Repository interface {
	GetProject(context.Context, int64) (project.Project, error)
	Create(context.Context, sprint.Sprint, project.Project) (sprint.Sprint, error)
	Get(context.Context, int64, int64) (sprint.Sprint, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) Service { return Service{repository: repository} }

func (s Service) Create(ctx context.Context, projectID int64, goal string) (sprint.Sprint, error) {
	value, err := sprint.New(projectID, goal)
	if err != nil {
		return sprint.Sprint{}, err
	}

	projectValue, err := s.repository.GetProject(ctx, projectID)
	if err != nil {
		return sprint.Sprint{}, err
	}
	projectValue.StartFirstSprint()
	return s.repository.Create(ctx, value, projectValue)
}

func (s Service) Get(ctx context.Context, projectID, sprintID int64) (sprint.Sprint, error) {
	if projectID < 1 || sprintID < 1 {
		return sprint.Sprint{}, sprint.ErrNotFound
	}
	return s.repository.Get(ctx, projectID, sprintID)
}
