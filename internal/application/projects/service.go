// Package projects coordinates project creation and consultation use cases.
package projects

import (
	"context"
	"errors"
	"fmt"
	"time"
	_ "time/tzdata" // PT-01: IANA zone is available on hosts without system zone files.

	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/google/uuid"
)

// Repository stores and retrieves the shared project entity.
type Repository interface {
	Save(context.Context, project.Project) error
	List(context.Context) ([]project.Project, error)
	FindByID(context.Context, string) (project.Project, error)
	// Mutate calls a repeatable, side-effect-free domain transformation after
	// locking and reading the current record; a confirmed rollback may retry it.
	Mutate(context.Context, string, func(project.Project) (project.Project, error)) (project.Project, error)
}

// List retrieves all saved projects for the home screen.
func (s *Service) List(ctx context.Context) ([]project.Project, error) {
	values, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return values, nil
}

// Service implements project creation and consultation use cases.
type Service struct {
	repository Repository
	clock      func() time.Time
}

// NewService constructs a project service backed by repository.
func NewService(repository Repository) *Service {
	return NewServiceWithClock(repository, time.Now)
}

// ErrBusy identifies bounded write contention without exposing driver errors.
var ErrBusy = errors.New("project is busy; retry later")

// NewServiceWithClock supports deterministic days without using the client clock.
func NewServiceWithClock(repository Repository, clock func() time.Time) *Service {
	return &Service{repository: repository, clock: clock}
}

func (s *Service) Update(ctx context.Context, id string, change project.Patch) (project.Project, error) {
	if err := change.Validate(); err != nil {
		return project.Project{}, err
	}
	return s.repository.Mutate(ctx, id, func(current project.Project) (project.Project, error) { return current.Edit(change) })
}

func (s *Service) Finalize(ctx context.Context, id string) (project.Project, error) {
	return s.repository.Mutate(ctx, id, func(current project.Project) (project.Project, error) {
		if current.Status == project.Finalized {
			return current, nil
		}
		zone, err := time.LoadLocation("America/Argentina/Buenos_Aires")
		if err != nil {
			return project.Project{}, fmt.Errorf("load project date zone: %w", err)
		}
		today, err := project.ParseDate(s.clock().In(zone).Format(time.DateOnly))
		if err != nil {
			return project.Project{}, err
		}
		return current.Finalize(today)
	})
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
