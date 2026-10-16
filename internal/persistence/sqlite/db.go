package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
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

	// Driver pragmas run for every new pooled connection, not just the initial one.
	dsn := ":memory:?_pragma=foreign_keys(1)"
	if path != ":memory:" {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("resolve SQLite path: %w", err)
		}
		uriPath := filepath.ToSlash(absolute)
		if filepath.VolumeName(absolute) != "" {
			uriPath = "/" + uriPath
		}
		dsn = (&url.URL{Scheme: "file", Path: uriPath, RawQuery: "_pragma=foreign_keys(1)"}).String()
	}
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open SQLite database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to SQLite database: %w", err)
	}

	return db, nil
}
