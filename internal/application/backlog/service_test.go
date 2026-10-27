package backlog

import (
	"context"
	"errors"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"reflect"
	"testing"

	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/google/uuid"
)

func TestServiceCreatesBacklogItemWithServerGeneratedUUID(t *testing.T) {
	repository := newMemoryRepository()
	service := NewService(repository)

	created, err := service.Create(context.Background(), "project-1", validInput())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := uuid.Parse(created.ID); err != nil {
		t.Errorf("Create() ID = %q, want server-generated UUID: %v", created.ID, err)
	}
	if created.ProjectID != "project-1" || created.Status != domain.Pending {
		t.Errorf("Create() = %+v, want project-1 and Pending", created)
	}
	if !reflect.DeepEqual(repository.items[created.ID], created) {
		t.Errorf("Create() did not persist item: %+v", created)
	}
}

func TestServiceRejectsMissingProject(t *testing.T) {
	service := NewService(newMemoryRepository())
	if _, err := service.Create(context.Background(), "missing", validInput()); !errors.Is(err, domain.ErrProjectNotFound) {
		t.Errorf("Create() error = %v, want ErrProjectNotFound", err)
	}
}

func TestServiceUpdatesOnlyEditableFields(t *testing.T) {
	repository := newMemoryRepository()
	service := NewService(repository)
	created, err := service.Create(context.Background(), "project-1", validInput())
	if err != nil {
		t.Fatal(err)
	}
	created.Status = domain.InProgress
	repository.items[created.ID] = created

	input := validInput()
	input.Title = "Updated"
	updated, err := service.Update(context.Background(), "project-1", created.ID, input)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Title != "Updated" || updated.Status != domain.InProgress || updated.ID != created.ID {
		t.Errorf("Update() = %+v, want edited title and preserved identity/status", updated)
	}
}

func TestServiceListsOnlyItemsForExistingProject(t *testing.T) {
	repository := newMemoryRepository()
	service := NewService(repository)
	created, err := service.Create(context.Background(), "project-1", validInput())
	if err != nil {
		t.Fatal(err)
	}

	items, err := service.List(context.Background(), "project-1")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 || items[0].ID != created.ID {
		t.Errorf("List() = %+v, want only item %q", items, created.ID)
	}
	if _, err := service.List(context.Background(), "missing"); !errors.Is(err, domain.ErrProjectNotFound) {
		t.Errorf("List(missing) error = %v, want ErrProjectNotFound", err)
	}
}

func validInput() domain.Input {
	return domain.Input{EstimatedHours: hours.Ptr(800),
		Title: "Login", Description: "User can sign in", Priority: domain.PriorityMedium,
		AcceptanceCriteria: []string{"Valid credentials grant access"},
	}
}

type memoryRepository struct {
	projects map[string]bool
	items    map[string]domain.Item
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{
		projects: map[string]bool{"project-1": true},
		items:    make(map[string]domain.Item),
	}
}

func (r *memoryRepository) ProjectExists(_ context.Context, id string) (bool, error) {
	return r.projects[id], nil
}

func (r *memoryRepository) Create(_ context.Context, item domain.Item) error {
	r.items[item.ID] = item
	return nil
}

func (r *memoryRepository) FindByID(_ context.Context, projectID, id string) (domain.Item, error) {
	item, ok := r.items[id]
	if !ok || item.ProjectID != projectID {
		return domain.Item{}, domain.ErrNotFound
	}
	return item, nil
}

func (r *memoryRepository) Update(_ context.Context, item domain.Item) error {
	if _, ok := r.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	r.items[item.ID] = item
	return nil
}

func (r *memoryRepository) ListByProject(_ context.Context, projectID string) ([]domain.Item, error) {
	var result []domain.Item
	for _, item := range r.items {
		if item.ProjectID == projectID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *memoryRepository) Mutate(ctx context.Context, pid, id string, change func(domain.Item) (domain.Item, error)) (domain.Item, error) {
	if !r.projects[pid] {
		return domain.Item{}, domain.ErrProjectNotFound
	}
	current, err := r.FindByID(ctx, pid, id)
	if err != nil {
		return current, err
	}
	next, err := change(current)
	if err != nil {
		return next, err
	}
	return next, r.Update(ctx, next)
}
