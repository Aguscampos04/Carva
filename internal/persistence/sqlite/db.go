package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	DatabasePathEnv     = "CARVA_DB_PATH"
	DefaultDatabasePath = ".local/carva.db"
	driverName          = "sqlite"
)

// DatabasePath returns the configured database path or the local default.
func DatabasePath() string {
	if path := strings.TrimSpace(os.Getenv(DatabasePathEnv)); path != "" {
		return path
	}

	return DefaultDatabasePath
}

// Open opens and verifies the database at the configured path.
func Open(ctx context.Context) (*sql.DB, error) {
	return OpenPath(ctx, DatabasePath())
}

// OpenPath opens and verifies a SQLite database at path.
func OpenPath(ctx context.Context, path string) (*sql.DB, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New("SQLite database path must not be empty")
	}

	if path != ":memory:" {
		directory := filepath.Dir(path)
		if directory != "." {
			if err := os.MkdirAll(directory, 0o755); err != nil {
				return nil, fmt.Errorf("create SQLite database directory: %w", err)
			}
		}
	}

	db, err := sql.Open(driverName, path)
	if err != nil {
		return nil, fmt.Errorf("open SQLite database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to SQLite database: %w", err)
	}

	return db, nil
}
