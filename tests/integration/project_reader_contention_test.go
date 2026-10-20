package integration_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

// PT-03 / BDD-25: a reader can allow BEGIN IMMEDIATE but temporarily prevent
// COMMIT. Retry must roll back and restart the whole mutation, within its bound.
func TestProjectFinalizationRetriesAfterReaderBlocksCommit(t *testing.T) {
	path, db := openMigratedProjectDB(t)
	other, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	entered := make(chan struct{})
	var once sync.Once
	service := projects.NewServiceWithClock(sqlite.NewProjectRepository(db), func() time.Time {
		once.Do(func() { close(entered) })
		return time.Date(2026, 10, 20, 15, 0, 0, 0, time.UTC)
	})
	p, err := service.Create(context.Background(), "Reader contention")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := other.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, err = conn.ExecContext(context.Background(), "BEGIN"); err != nil {
		t.Fatal(err)
	}
	defer conn.ExecContext(context.Background(), "ROLLBACK")
	var count int
	if err = conn.QueryRowContext(context.Background(), `SELECT count(*) FROM projects`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, err := service.Finalize(context.Background(), p.ID); done <- err }()
	<-entered
	// Keep the shared read lock until at least one attempted commit has failed.
	time.Sleep(50 * time.Millisecond)
	if _, err = conn.ExecContext(context.Background(), "ROLLBACK"); err != nil {
		t.Fatal(err)
	}
	if err = <-done; err != nil {
		t.Fatalf("expected successful bounded retry, got %v", err)
	}
	got, err := service.Get(context.Background(), p.ID)
	if err != nil || got.EndDate != "2026-10-20" {
		t.Fatalf("%+v %v", got, err)
	}
}
