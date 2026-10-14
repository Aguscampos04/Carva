package integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	carvasqlite "github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

func TestProjectCreationAndLookupPersistAcrossDatabaseReopen(t *testing.T) {
	ctx := context.Background()
	databasePath, db := openMigratedProjectDB(t)
	service := projects.NewService(carvasqlite.NewProjectRepository(db))

	first, err := service.Create(ctx, "Carva MVP")
	if err != nil {
		t.Fatalf("Create(first) error = %v", err)
	}
	second, err := service.Create(ctx, "Proyecto de prueba")
	if err != nil {
		t.Fatalf("Create(second) error = %v", err)
	}
	if first.ID == second.ID {
		t.Fatalf("created projects share ID %q", first.ID)
	}
	if first.Status != project.Planned || second.Status != project.Planned {
		t.Fatalf("created states = %q and %q, want %q", first.Status, second.Status, project.Planned)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close() before reopen: %v", err)
	}

	reopened, err := carvasqlite.OpenPath(ctx, databasePath)
	if err != nil {
		t.Fatalf("reopen database: %v", err)
	}
	t.Cleanup(func() { _ = reopened.Close() })
	reader := projects.NewService(carvasqlite.NewProjectRepository(reopened))
	for _, want := range []project.Project{first, second} {
		got, err := reader.Get(ctx, want.ID)
		if err != nil {
			t.Errorf("Get(%q) error = %v", want.ID, err)
			continue
		}
		if got != want {
			t.Errorf("Get(%q) = %+v, want %+v", want.ID, got, want)
		}
	}
}

func TestProjectLookupReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	_, db := openMigratedProjectDB(t)
	service := projects.NewService(carvasqlite.NewProjectRepository(db))
	if _, err := service.Get(ctx, "missing"); !errors.Is(err, project.ErrNotFound) {
		t.Errorf("Get(missing) error = %v, want %v", err, project.ErrNotFound)
	}
}
