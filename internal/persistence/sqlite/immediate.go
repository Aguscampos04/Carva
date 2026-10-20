package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/Aguscampos04/Carva/internal/application/projects"
)

// withImmediate serializes project edits, finalization and Sprint creation across
// independent SQLite connections/processes. The callback runs only after locking.
func withImmediate(ctx context.Context, db *sql.DB, change func(*sql.Conn) error) error {
	waitCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := db.Conn(waitCtx)
	if err != nil {
		return contentionError(waitCtx, err)
	}
	defer conn.Close()
	for {
		retry, err := immediateAttempt(waitCtx, conn, change)
		if err == nil {
			return nil
		}
		if !retry {
			return contentionError(waitCtx, err)
		}
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case <-waitCtx.Done():
			timer.Stop()
			return fmt.Errorf("%w: %w", projects.ErrBusy, waitCtx.Err())
		case <-timer.C:
		}
	}
}

// immediateAttempt permits retry only after a definite BUSY/LOCKED failure and
// confirmed rollback. Arbitrary commit errors have an uncertain outcome and
// must be returned to the client, never replayed automatically.
func immediateAttempt(ctx context.Context, conn *sql.Conn, change func(*sql.Conn) error) (retry bool, err error) {
	if _, err = conn.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
		return sqliteBusy(err), err
	}
	committed := false
	defer func() {
		if !committed {
			cleanup, stop := context.WithTimeout(context.Background(), time.Second)
			defer stop()
			if _, rollbackErr := conn.ExecContext(cleanup, "ROLLBACK"); rollbackErr != nil {
				// Never return a connection with a live manual transaction to the pool.
				_ = conn.Raw(func(any) error { return driver.ErrBadConn })
				retry = false
				err = errors.Join(err, rollbackErr)
			}
		}
	}()
	if err := change(conn); err != nil {
		return sqliteBusy(err), err
	}
	if _, err := conn.ExecContext(ctx, "COMMIT"); err != nil {
		return sqliteBusy(err), err
	}
	committed = true
	return false, nil
}

func sqliteBusy(err error) bool {
	var coded interface{ Code() int }
	return errors.As(err, &coded) && (coded.Code()&255 == 5 || coded.Code()&255 == 6)
}
func contentionError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return fmt.Errorf("%w: %w", projects.ErrBusy, ctx.Err())
	}
	return err
}
