package sprints

import (
	"context"
	"time"

	"github.com/Aguscampos04/Carva/internal/domain/backlog"
	effortdomain "github.com/Aguscampos04/Carva/internal/domain/effort"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
	"github.com/google/uuid"
	_ "time/tzdata"
)

type Repository interface {
	GetProject(context.Context, string) (project.Project, error)
	Create(context.Context, sprint.Sprint, func(project.Project) (project.Project, error)) (sprint.Sprint, error)
	Get(context.Context, string, int64) (sprint.Sprint, error)
	ListClosed(context.Context, string) ([]sprint.Sprint, error)
	Close(context.Context, string, int64) (sprint.Sprint, error)
	AssignBacklogItem(context.Context, string, int64, sprint.PlannedBacklogItem) error
	ListPlannedBacklogItems(context.Context, string, int64) ([]sprint.PlannedBacklogItem, error)
	CompleteBacklogItem(context.Context, string, int64, string, bool) (sprint.PlannedBacklogItem, error)
	CompleteBacklogItemWithEffort(context.Context, string, int64, string, bool, *effortdomain.Entry, func(project.Project, hours.Value) error) (sprint.PlannedBacklogItem, error)
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

	return s.repository.Create(ctx, value, func(current project.Project) (project.Project, error) {
		if err := current.CanCreateSprint(); err != nil {
			return project.Project{}, err
		}
		current.StartFirstSprint()
		return current, nil
	})
}

func (s Service) Get(ctx context.Context, projectID string, sprintID int64) (sprint.Sprint, error) {
	if sprintID < 1 || projectID == "" {
		return sprint.Sprint{}, sprint.ErrNotFound
	}
	return s.repository.Get(ctx, projectID, sprintID)
}

// ListClosed returns only the persisted closed Sprints of an existing project.
func (s Service) ListClosed(ctx context.Context, projectID string) ([]sprint.Sprint, error) {
	if projectID == "" {
		return nil, project.ErrNotFound
	}
	if _, err := s.repository.GetProject(ctx, projectID); err != nil {
		return nil, err
	}
	return s.repository.ListClosed(ctx, projectID)
}

// Close closes a Sprint and rolls its incomplete stories into the next available Sprint.
func (s Service) Close(ctx context.Context, projectID string, sprintID int64) (sprint.Sprint, error) {
	if projectID == "" || sprintID < 1 {
		return sprint.Sprint{}, sprint.ErrNotFound
	}
	return s.repository.Close(ctx, projectID, sprintID)
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
	values, err := s.repository.ListPlannedBacklogItems(ctx, projectID, sprintID)
	if err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	for _, v := range values {
		if v.BacklogItemID == itemID {
			return v, nil
		}
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

func (s Service) CompleteBacklogItem(ctx context.Context, projectID string, sprintID int64, itemID string, verified bool) (sprint.PlannedBacklogItem, error) {
	return s.completeBacklogItem(ctx, projectID, sprintID, itemID, verified, nil)
}

// CompleteBacklogItemWithEffort records a required new effort entry and completes
// the story atomically when it has no previously active effort.
func (s Service) CompleteBacklogItemWithEffort(ctx context.Context, projectID string, sprintID int64, itemID string, verified bool, effortInput *effortdomain.Input) (sprint.PlannedBacklogItem, error) {
	var entry *effortdomain.Entry
	var check func(project.Project, hours.Value) error
	if effortInput != nil {
		effortInput.StoryID = itemID
		effortInput.SprintID = &sprintID
		effortInput.HasSprint = true
		if err := effortInput.Validate(); err != nil {
			return sprint.PlannedBacklogItem{}, err
		}
		entry = &effortdomain.Entry{ID: uuid.NewString(), ProjectID: projectID, Input: *effortInput}
		check = func(value project.Project, other hours.Value) error {
			location, err := time.LoadLocation("America/Argentina/Buenos_Aires")
			if err != nil {
				return err
			}
			today := project.Date(time.Now().In(location).Format(time.DateOnly))
			if err := effortdomain.ValidateDate(value, today, effortInput.WorkDate); err != nil {
				return err
			}
			return effortdomain.ValidateDaily(other, effortInput.Hours)
		}
	}
	return s.completeBacklogItem(ctx, projectID, sprintID, itemID, verified, entry, check)
}

func (s Service) completeBacklogItem(ctx context.Context, projectID string, sprintID int64, itemID string, verified bool, entry *effortdomain.Entry, validateEffort ...func(project.Project, hours.Value) error) (sprint.PlannedBacklogItem, error) {
	if projectID == "" || sprintID < 1 {
		return sprint.PlannedBacklogItem{}, sprint.ErrNotFound
	}
	if itemID == "" {
		return sprint.PlannedBacklogItem{}, sprint.ErrInvalidBacklogItemID
	}
	if !verified {
		return sprint.PlannedBacklogItem{}, backlog.ErrAcceptanceCriteriaUnverified
	}
	if _, err := s.repository.Get(ctx, projectID, sprintID); err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	item, err := s.backlogRepository.FindByID(ctx, projectID, itemID)
	if err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	if _, err := item.Complete(verified); err != nil {
		return sprint.PlannedBacklogItem{}, err
	}
	var check func(project.Project, hours.Value) error
	if len(validateEffort) > 0 {
		check = validateEffort[0]
	}
	return s.repository.CompleteBacklogItemWithEffort(ctx, projectID, sprintID, itemID, verified, entry, check)
}
