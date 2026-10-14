package integration_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	carvasqlite "github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

func openMigratedProjectDB(t *testing.T) (string, *sql.DB) {
	t.Helper()
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "project.db")
	db, err := carvasqlite.OpenPath(ctx, databasePath)
	if err != nil {
		t.Fatalf("OpenPath() error = %v", err)
	}
	if err := carvasqlite.Migrate(ctx, db); err != nil {
		_ = db.Close()
		t.Fatalf("Migrate() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return databasePath, db
}
