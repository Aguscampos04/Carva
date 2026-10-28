package projects

import (
	"context"
	"errors"
	"testing"

	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/google/uuid"
)

func TestServiceCreatesAndRetrievesProject(t *testing.T) {
	store := newMemoryRepository()
	service := NewService(store)

	created, err := service.Create(context.Background(), "Carva MVP")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := uuid.Parse(created.ID); err != nil {
		t.Errorf("Create() ID = %q, want UUID: %v", created.ID, err)
	}
	if created.Name != "Carva MVP" || created.Status != project.Planned {
		t.Errorf("Create() = %+v, want entered name and Planned status", created)
	}

	got, err := service.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != created {
		t.Errorf("Get() = %+v, want %+v", got, created)
	}
}

func TestServiceRejectsInvalidNameBeforePersistence(t *testing.T) {
	store := newMemoryRepository()
	service := NewService(store)

	_, err := service.Create(context.Background(), " \t ")
	if !errors.Is(err, project.ErrInvalidName) {
		t.Errorf("Create() error = %v, want ErrInvalidName", err)
	}
	if len(store.projects) != 0 {
		t.Errorf("invalid create persisted %d projects, want none", len(store.projects))
	}
}

func TestServiceReturnsRepositoryFailure(t *testing.T) {
	wantErr := errors.New("storage unavailable")
	service := NewService(&memoryRepository{saveErr: wantErr})

	_, err := service.Create(context.Background(), "Carva MVP")
	if !errors.Is(err, wantErr) {
		t.Errorf("Create() error = %v, want wrapped storage error", err)
	}
}

func TestServiceReturnsNotFound(t *testing.T) {
	service := NewService(newMemoryRepository())

	_, err := service.Get(context.Background(), "missing")
	if !errors.Is(err, project.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

type memoryRepository struct {
	projects map[string]project.Project
	saveErr  error
}

func (r *memoryRepository) List(_ context.Context) ([]project.Project, error) {
	values := make([]project.Project, 0, len(r.projects))
	for _, value := range r.projects {
		values = append(values, value)
	}
	return values, r.saveErr
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{projects: make(map[string]project.Project)}
}

func (r *memoryRepository) Save(_ context.Context, p project.Project) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.projects[p.ID] = p
	return nil
}

func (r *memoryRepository) FindByID(_ context.Context, id string) (project.Project, error) {
	p, ok := r.projects[id]
	if !ok {
		return project.Project{}, project.ErrNotFound
	}
	return p, nil
}

func (r *memoryRepository) Mutate(ctx context.Context, id string, change func(project.Project) (project.Project, error)) (project.Project, error) {
	p, err := r.FindByID(ctx, id)
	if err != nil {
		return project.Project{}, err
	}
	next, err := change(p)
	if err != nil {
		return project.Project{}, err
	}
	if r.saveErr != nil {
		return project.Project{}, r.saveErr
	}
	r.projects[id] = next
	return next, nil
}
