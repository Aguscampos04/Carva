package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Aguscampos04/Carva/internal/persistence/sqlite/migrations"
)

const migrationTable = "carva_schema_migrations"

// Migrate applies embedded, versioned SQL migrations that have not yet run.
func Migrate(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("SQLite database must not be nil")
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS `+migrationTable+` (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create SQLite migration ledger: %w", err)
	}

	entries, err := fs.ReadDir(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("read embedded SQLite migrations: %w", err)
	}

	seenVersions := make(map[int64]string)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".sql" {
			continue
		}

		version, err := migrationVersion(name)
		if err != nil {
			return err
		}
		if previous, exists := seenVersions[version]; exists {
			return fmt.Errorf("SQLite migrations %q and %q use the same version", previous, name)
		}
		seenVersions[version] = name

		var applied bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM `+migrationTable+` WHERE version = ?)`, version,
		).Scan(&applied); err != nil {
			return fmt.Errorf("check SQLite migration %q: %w", name, err)
		}
		if applied {
			continue
		}

		sqlText, err := fs.ReadFile(migrations.Files, name)
		if err != nil {
			return fmt.Errorf("read SQLite migration %q: %w", name, err)
		}

		if err := applyMigration(ctx, db, version, name, string(sqlText)); err != nil {
			return err
		}
	}

	return nil
}

func migrationVersion(name string) (int64, error) {
	prefix, _, ok := strings.Cut(name, "_")
	if !ok || len(prefix) != 6 {
		return 0, fmt.Errorf("SQLite migration %q must start with a six-digit version and underscore", name)
	}

	version, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil || version < 1 {
		return 0, fmt.Errorf("SQLite migration %q has an invalid version", name)
	}

	return version, nil
}

func applyMigration(ctx context.Context, db *sql.DB, version int64, name, sqlText string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin SQLite migration %q: %w", name, err)
	}

	if _, err := tx.ExecContext(ctx, sqlText); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("execute SQLite migration %q: %w", name, err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO `+migrationTable+` (version, name) VALUES (?, ?)`, version, name,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record SQLite migration %q: %w", name, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit SQLite migration %q: %w", name, err)
	}

	return nil
}
