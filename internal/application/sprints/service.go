package sprints

import (
	"context"

	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type Repository interface {
	GetProject(context.Context, string) (project.Project, error)
	Create(context.Context, sprint.Sprint, project.Project) (sprint.Sprint, error)
	Get(context.Context, string, int64) (sprint.Sprint, error)
	AssignBacklogItem(context.Context, string, int64, sprint.PlannedBacklogItem) error
	ListPlannedBacklogItems(context.Context, string, int64) ([]sprint.PlannedBacklogItem, error)
}

// BacklogRepository supplies the current item values to snapshot at assignment time.
type BacklogRepository interface {
	FindByID(context.Context, string, string) (backlog.Item, error)
}

type Service struct {
	repository        Repository
	backlogRepository BacklogRepository
}

func NewService(repository Repository, backlogRepository BacklogRepository) Service {
	return Service{repository: repository, backlogRepository: backlogRepository}
}

func (s Service) Create(ctx context.Context, projectID string, goal string) (sprint.Sprint, error) {
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

func (s Service) Get(ctx context.Context, projectID string, sprintID int64) (sprint.Sprint, error) {
	if sprintID < 1 || projectID == "" {
		return sprint.Sprint{}, sprint.ErrNotFound
	}
	return s.repository.Get(ctx, projectID, sprintID)
}

func (s Service) AssignBacklogItem(ctx context.Context, projectID string, sprintID int64, itemID string) (sprint.PlannedBacklogItem, error) {
	if projectID == "" || sprintID < 1 {
		return sprint.PlannedBacklogItem{}, sprint.ErrNotFound
	}
	if itemID == "" {
		return sprint.PlannedBacklogItem{}, sprint.ErrInvalidBacklogItemID
	}
	if _, err := s.repository.Get(ctx, projectID, sprintID); err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	item, err := s.backlogRepository.FindByID(ctx, projectID, itemID)
	if err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	planned, err := sprint.Snapshot(sprintID, item)
	if err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	if err := s.repository.AssignBacklogItem(ctx, projectID, sprintID, planned); err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	return planned, nil
}

func (s Service) ListPlannedBacklogItems(ctx context.Context, projectID string, sprintID int64) ([]sprint.PlannedBacklogItem, error) {
	if sprintID < 1 || projectID == "" {
		return nil, sprint.ErrNotFound
	}
	if _, err := s.repository.Get(ctx, projectID, sprintID); err != nil {
		return nil, err
	}
	return s.repository.ListPlannedBacklogItems(ctx, projectID, sprintID)
}
