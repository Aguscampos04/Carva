package integration_test

import (
	"context"
	"errors"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"reflect"
	"testing"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

func TestBacklogPersistsCreatesUpdatesAndListsProjectItems(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projectRepository := sqlite.NewProjectRepository(db)
	firstProject, err := project.New("project-1", "Project One")
	if err != nil {
		t.Fatal(err)
	}
	if err := projectRepository.Save(context.Background(), firstProject); err != nil {
		t.Fatal(err)
	}
	service := appbacklog.NewService(sqlite.NewBacklogRepository(db))

	provisional := 3
	input := backlogInput()
	input.ProvisionalStoryPoints = &provisional
	created, err := service.Create(context.Background(), firstProject.ID, input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Status != domain.Pending || created.ProvisionalStoryPoints == nil || *created.ProvisionalStoryPoints != provisional {
		t.Fatalf("Create() = %+v, want Pending with provisional points %d", created, provisional)
	}

	updatedInput := backlogInput()
	updatedInput.Title = "Updated title"
	updated, err := service.Update(context.Background(), firstProject.ID, created.ID, updatedInput)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Title != updatedInput.Title || updated.ProvisionalStoryPoints != nil {
		t.Errorf("Update() = %+v, want updated fields", updated)
	}

	items, err := service.List(context.Background(), firstProject.ID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 || !reflect.DeepEqual(items[0], updated) {
		t.Errorf("List() = %+v, want persisted item %+v", items, updated)
	}
}

func TestBacklogListIsProjectScopedAndAllowsEmptyProject(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projects := sqlite.NewProjectRepository(db)
	for _, p := range []project.Project{
		{ID: "project-1", Name: "One", Status: project.Planned},
		{ID: "project-2", Name: "Two", Status: project.Planned},
	} {
		if err := projects.Save(context.Background(), p); err != nil {
			t.Fatal(err)
		}
	}
	service := appbacklog.NewService(sqlite.NewBacklogRepository(db))
	if _, err := service.Create(context.Background(), "project-1", backlogInput()); err != nil {
		t.Fatal(err)
	}
	if got, err := service.List(context.Background(), "project-2"); err != nil || len(got) != 0 {
		t.Errorf("List(project-2) = (%+v, %v), want empty list", got, err)
	}
	if _, err := service.List(context.Background(), "missing"); !errors.Is(err, domain.ErrProjectNotFound) {
		t.Errorf("List(missing) error = %v, want ErrProjectNotFound", err)
	}
}

func TestBacklogRejectsUpdatingAnItemThroughAnotherProject(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projects := sqlite.NewProjectRepository(db)
	for _, p := range []project.Project{
		{ID: "project-1", Name: "One", Status: project.Planned},
		{ID: "project-2", Name: "Two", Status: project.Planned},
	} {
		if err := projects.Save(context.Background(), p); err != nil {
			t.Fatal(err)
		}
	}
	service := appbacklog.NewService(sqlite.NewBacklogRepository(db))
	created, err := service.Create(context.Background(), "project-1", backlogInput())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Update(context.Background(), "project-2", created.ID, backlogInput()); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Update(other project) error = %v, want ErrNotFound", err)
	}
}

func backlogInput() domain.Input {
	return domain.Input{EstimatedHours: hours.Ptr(800),
		Title: "Login", Description: "User can sign in", Priority: domain.PriorityHigh,
		AcceptanceCriteria: []string{"Valid credentials grant access"},
	}
}
