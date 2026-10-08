package integration_test

import (
	"context"
	"path/filepath"
	"testing"

	carvasqlite "github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

func TestSQLiteConnectionAndMigrationsUseTemporaryDatabase(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "nested", "integration.db")

	db, err := carvasqlite.OpenPath(ctx, databasePath)
	if err != nil {
		t.Fatalf("OpenPath() error = %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	if err := carvasqlite.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	if err := carvasqlite.Migrate(ctx, db); err != nil {
		t.Fatalf("second Migrate() error = %v", err)
	}

	var ledgerExists bool
	if err := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'carva_schema_migrations')`,
	).Scan(&ledgerExists); err != nil {
		t.Fatalf("query migration ledger: %v", err)
	}
	if !ledgerExists {
		t.Fatal("migration ledger table was not created")
	}
	var recordedMigrations int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM carva_schema_migrations`).Scan(&recordedMigrations); err != nil {
		t.Fatalf("query applied migrations: %v", err)
	}
	if recordedMigrations != 1 {
		t.Fatalf("expected the baseline migration to be recorded once, got %d", recordedMigrations)
	}

	var businessTables int
	if err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name != 'carva_schema_migrations'`,
	).Scan(&businessTables); err != nil {
		t.Fatalf("query application tables: %v", err)
	}
	if businessTables != 0 {
		t.Fatalf("expected no application tables, got %d", businessTables)
	}
}

func TestSQLiteTemporaryDatabasesAreIsolated(t *testing.T) {
	ctx := context.Background()
	first, err := carvasqlite.OpenPath(ctx, filepath.Join(t.TempDir(), "first.db"))
	if err != nil {
		t.Fatalf("open first temporary database: %v", err)
	}
	t.Cleanup(func() { _ = first.Close() })

	second, err := carvasqlite.OpenPath(ctx, filepath.Join(t.TempDir(), "second.db"))
	if err != nil {
		t.Fatalf("open second temporary database: %v", err)
	}
	t.Cleanup(func() { _ = second.Close() })

	if _, err := first.ExecContext(ctx, `CREATE TABLE integration_probe (id INTEGER PRIMARY KEY)`); err != nil {
		t.Fatalf("create table in first database: %v", err)
	}

	var tableCount int
	if err := second.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'integration_probe'`,
	).Scan(&tableCount); err != nil {
		t.Fatalf("check second database: %v", err)
	}
	if tableCount != 0 {
		t.Fatalf("expected second database not to contain first database's table, got %d", tableCount)
	}
}

func TestDatabasePathUsesEnvironmentAndLocalDefault(t *testing.T) {
	t.Setenv(carvasqlite.DatabasePathEnv, "")
	if got := carvasqlite.DatabasePath(); got != carvasqlite.DefaultDatabasePath {
		t.Fatalf("DatabasePath() default = %q, want %q", got, carvasqlite.DefaultDatabasePath)
	}

	customPath := filepath.Join(t.TempDir(), "custom.db")
	t.Setenv(carvasqlite.DatabasePathEnv, customPath)
	if got := carvasqlite.DatabasePath(); got != customPath {
		t.Fatalf("DatabasePath() = %q, want %q", got, customPath)
	}
}
