package integration_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite/migrations"
)

// BDD-21/22/27/28: independent DB pools exercise the real SQLite boundary.
func TestProjectLifecycleConcurrentWriters(t *testing.T) {
	path, db := openMigratedProjectDB(t)
	other, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	handlers := []http.Handler{lifecycleHandler(db, "2026-10-20T15:00:00Z"), lifecycleHandler(other, "2026-10-20T15:00:00Z")}
	for round := 0; round < 8; round++ {
		p := lifeRequest(t, handlers[0], "POST", "/api/projects", `{"name":"Demo"}`, 201)
		url := "/api/projects/" + p["id"].(string)
		operations := []struct {
			method, suffix, body string
			allowed              []int
		}{
			{"PATCH", "", `{"name":"Concurrent"}`, []int{200}},
			{"PATCH", "", `{"start_date":"2026-10-01"}`, []int{200}},
			{"POST", "/finalize", `{"confirmed":true}`, []int{200}},
			{"POST", "/finalize", `{"confirmed":true}`, []int{200}},
			{"POST", "/sprints", `{"goal":"Concurrent Sprint"}`, []int{201, 409}},
		}
		start := make(chan struct{})
		var wg sync.WaitGroup
		codes := make([]int, len(operations))
		for i, op := range operations {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				<-start
				w := httptest.NewRecorder()
				handlers[i%2].ServeHTTP(w, httptest.NewRequest(op.method, url+op.suffix, strings.NewReader(op.body)))
				codes[i] = w.Code
				valid := false
				for _, code := range op.allowed {
					valid = valid || code == w.Code
				}
				if !valid {
					t.Errorf("concurrent %s: %d %s", op.suffix, w.Code, w.Body)
				}
			}(i)
		}
		close(start)
		wg.Wait()
		got := lifeRequest(t, handlers[0], "GET", url, "", 200)
		if got["name"] != "Concurrent" || got["start_date"] != "2026-10-01" || got["end_date"] != "2026-10-20" || got["status"] != "Finalizado" {
			t.Fatal(got)
		}
		var count int
		_ = db.QueryRow(`SELECT count(*) FROM sprints WHERE project_id=?`, p["id"]).Scan(&count)
		expected := 0
		if codes[4] == 201 {
			expected = 1
		}
		if count != expected {
			t.Fatalf("Sprint result %d, count %d", codes[4], count)
		}
	}
}

// BDD-22: enforce the finalization-first serialization order. Acquiring a lock
// alone does not establish that order: a BUSY COMMIT can roll back and retry.
func TestProjectFinalizationWinsBeforeSprintCreation(t *testing.T) {
	path, db := openMigratedProjectDB(t)
	other, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	entered, release := make(chan struct{}), make(chan struct{})
	service := projects.NewServiceWithClock(sqlite.NewProjectRepository(db), gatedProjectClock(entered, release))
	p, err := service.Create(context.Background(), "Demo")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, err := service.Finalize(context.Background(), p.ID); done <- err }()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("finalization did not acquire lock")
	}
	response := make(chan *httptest.ResponseRecorder, 1)
	finalized := make(chan struct{})
	go func() {
		<-finalized
		w := performJSON(t, lifecycleHandler(other, "2026-10-20T15:00:00Z"), "POST", "/api/projects/"+p.ID+"/sprints", `{"goal":"Too late"}`)
		response <- w
	}()
	close(release)
	finalizationErr := <-done
	close(finalized)
	if finalizationErr != nil {
		t.Fatal(finalizationErr)
	}
	if w := <-response; w.Code != 409 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}

// This clock is called inside a transformation that may be replayed after BUSY.
func gatedProjectClock(entered, release chan struct{}) func() time.Time {
	var enteredOnce sync.Once
	return func() time.Time {
		// Signal the first lock acquisition only. A confirmed SQLite rollback
		// may replay Finalize's callback; later calls still use the same gate/day.
		enteredOnce.Do(func() { close(entered) })
		<-release
		return time.Date(2026, 10, 20, 15, 0, 0, 0, time.UTC)
	}
}

// A read transaction deterministically prevents the first writer COMMIT. The
// second callback invocation proves rollback/retry, rather than relying on timing.
func TestProjectFinalizationClockSurvivesCommitRetry(t *testing.T) {
	path, db := openMigratedProjectDB(t)
	other, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	entered, release := make(chan struct{}), make(chan struct{})
	retried, allowCommit := make(chan struct{}), make(chan struct{})
	var releaseOnce, commitOnce sync.Once
	defer releaseOnce.Do(func() { close(release) })
	defer commitOnce.Do(func() { close(allowCommit) })
	gate := gatedProjectClock(entered, release)
	calls := 0 // Finalize invokes its transformation sequentially, including retries.
	service := projects.NewServiceWithClock(sqlite.NewProjectRepository(db), func() time.Time {
		calls++
		if calls == 2 {
			close(retried)
			<-allowCommit
		}
		return gate()
	})
	p, err := service.Create(ctx, "Commit retry")
	if err != nil {
		t.Fatal(err)
	}
	reader, err := other.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Rollback()
	var name string
	if err := reader.QueryRowContext(ctx, `SELECT name FROM projects WHERE id=?`, p.ID).Scan(&name); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, err := service.Finalize(ctx, p.ID); done <- err }()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal("finalization did not acquire writer lock")
	}
	releaseOnce.Do(func() { close(release) })
	select {
	case <-retried:
	case <-ctx.Done():
		t.Fatal("reader did not force a COMMIT retry")
	}
	if err := reader.Rollback(); err != nil {
		t.Fatal(err)
	}
	commitOnce.Do(func() { close(allowCommit) })
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("finalization did not finish after releasing reader")
	}
	got, err := service.Get(ctx, p.ID)
	if err != nil || got.Status != project.Finalized || got.EndDate != project.Date("2026-10-20") {
		t.Fatalf("finalized project = %+v, %v", got, err)
	}
}

func TestProjectSprintCreationWinsBeforeFinalization(t *testing.T) {
	path, db := openMigratedProjectDB(t)
	other, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	service := projects.NewService(sqlite.NewProjectRepository(db))
	p, err := service.Create(context.Background(), "Demo")
	if err != nil {
		t.Fatal(err)
	}
	entered, release := make(chan struct{}), make(chan struct{})
	created := make(chan error, 1)
	value, err := sprint.New(p.ID, "First")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_, err := sqlite.NewSprintRepository(db).Create(context.Background(), value, func(current project.Project) (project.Project, error) {
			close(entered)
			<-release
			if err := current.CanCreateSprint(); err != nil {
				return project.Project{}, err
			}
			current.StartFirstSprint()
			return current, nil
		})
		created <- err
	}()
	<-entered
	response := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		response <- performJSON(t, lifecycleHandler(other, "2026-10-20T15:00:00Z"), "POST", "/api/projects/"+p.ID+"/finalize", `{"confirmed":true}`)
	}()
	close(release)
	if err := <-created; err != nil {
		t.Fatal(err)
	}
	if w := <-response; w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	got := lifeRequest(t, lifecycleHandler(db, "2026-10-20T15:00:00Z"), "GET", "/api/projects/"+p.ID+"/sprints/1", "", 200)
	if got["goal"] != "First" {
		t.Fatal("created Sprint not preserved")
	}
}

// BDD-27: both serialization orders reject one of the incompatible writes.
func TestProjectConcurrentFutureStartAndFinalization(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	h := lifecycleHandler(db, "2026-10-20T15:00:00Z")
	for _, editFirst := range []bool{true, false} {
		p := lifeRequest(t, h, "POST", "/api/projects", `{"name":"Demo"}`, 201)
		url := "/api/projects/" + p["id"].(string)
		repo := sqlite.NewProjectRepository(db)
		entered, release := make(chan struct{}), make(chan struct{})
		done := make(chan error, 1)
		go func() {
			_, err := repo.Mutate(context.Background(), p["id"].(string), func(current project.Project) (project.Project, error) {
				close(entered)
				<-release
				if editFirst {
					return current.Edit(project.Patch{HasStartDate: true, StartDate: project.Date("2026-10-21")})
				}
				return current.Finalize(project.Date("2026-10-20"))
			})
			done <- err
		}()
		<-entered
		response := make(chan *httptest.ResponseRecorder, 1)
		go func() {
			method, suffix, body := "PATCH", "", `{"start_date":"2026-10-21"}`
			if editFirst {
				method, suffix, body = "POST", "/finalize", `{"confirmed":true}`
			}
			response <- performJSON(t, h, method, url+suffix, body)
		}()
		close(release)
		if err := <-done; err != nil {
			t.Fatal(err)
		}
		w := <-response
		want := 400
		if editFirst {
			want = 409
		}
		if w.Code != want {
			t.Fatalf("%d %s", w.Code, w.Body)
		}
		got := lifeRequest(t, h, "GET", url, "", 200)
		if editFirst {
			if got["status"] != "Planificado" || got["end_date"] != nil {
				t.Fatal(got)
			}
		} else {
			if got["status"] != "Finalizado" || got["start_date"] != nil {
				t.Fatal(got)
			}
		}
	}
}

// BDD-24/25: rollback, lock deadline, context cancellation, and pool recovery.
func TestProjectLifecycleRollbackAndContention(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	h := lifecycleHandler(db, "2026-10-20T15:00:00Z")
	p := lifeRequest(t, h, "POST", "/api/projects", `{"name":"Demo"}`, 201)
	url := "/api/projects/" + p["id"].(string)
	for _, test := range []struct{ trigger, method, suffix, body string }{
		{`CREATE TRIGGER reject_update BEFORE UPDATE ON projects BEGIN SELECT RAISE(ABORT,'controlled failure'); END`, "PATCH", "", `{"name":"Changed","start_date":"2026-10-01"}`},
		{`CREATE TRIGGER reject_update BEFORE UPDATE ON projects BEGIN SELECT RAISE(ABORT,'controlled failure'); END`, "POST", "/finalize", `{"confirmed":true}`},
		{`CREATE TRIGGER reject_update BEFORE INSERT ON sprints BEGIN SELECT RAISE(ABORT,'controlled failure'); END`, "POST", "/sprints", `{"goal":"Fails after transition"}`},
	} {
		if _, err := db.Exec(test.trigger); err != nil {
			t.Fatal(err)
		}
		response := lifeRequest(t, h, test.method, url+test.suffix, test.body, 500)
		if strings.Contains(fmt.Sprint(response), "controlled failure") {
			t.Fatal("SQL leaked")
		}
		if got := lifeRequest(t, h, "GET", url, "", 200); !reflect.DeepEqual(p, got) {
			t.Fatalf("partial: %v", got)
		}
		if _, err := db.Exec(`DROP TRIGGER reject_update`); err != nil {
			t.Fatal(err)
		}
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = conn.ExecContext(context.Background(), "BEGIN IMMEDIATE"); err != nil {
		t.Fatal(err)
	}
	// The public API must stop waiting after PT-03's five-second bound.
	started := time.Now()
	lifeRequest(t, h, "POST", url+"/finalize", `{"confirmed":true}`, 503)
	if elapsed := time.Since(started); elapsed < 4*time.Second || elapsed > 8*time.Second {
		t.Fatalf("write bound: %s", elapsed)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	_, err = projects.NewService(sqlite.NewProjectRepository(db)).Finalize(ctx, p["id"].(string))
	if !errors.Is(err, projects.ErrBusy) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline: %v", err)
	}
	if _, err = conn.ExecContext(context.Background(), "ROLLBACK"); err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()
	lifeRequest(t, h, "POST", url+"/finalize", `{"confirmed":true}`, 200)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	lifeRequest(t, h, "PATCH", url, `{"name":"Unstored"}`, 500)
	lifeRequest(t, h, "POST", url+"/finalize", `{"confirmed":true}`, 500)
}

func TestProjectDatesMigrationPreservesExistingData(t *testing.T) {
	db, err := sqlite.OpenPath(context.Background(), filepath.Join(t.TempDir(), "old.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE carva_schema_migrations(version INTEGER PRIMARY KEY,name TEXT NOT NULL,applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		t.Fatal(err)
	}
	for i, name := range []string{"000001_baseline.sql", "000002_projects.sql", "000003_backlog_items.sql", "000004_sprints.sql", "000005_sprint_backlog_items.sql", "000006_members.sql", "000007_members_name_storage.sql"} {
		data, err := migrations.Files.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = db.Exec(string(data)); err != nil {
			t.Fatal(err)
		}
		if _, err = db.Exec(`INSERT INTO carva_schema_migrations(version,name) VALUES(?,?)`, i+1, name); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = db.Exec(`INSERT INTO projects VALUES('old','Existing','En progreso'),('legacy','Historical','Finalizado'); INSERT INTO sprints(id,project_id,goal) VALUES(1,'old','MVP'); INSERT INTO members VALUES('member','old','Ana')`); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err = sqlite.Migrate(context.Background(), db); err != nil {
			t.Fatal(err)
		}
	}
	p, err := sqlite.NewProjectRepository(db).FindByID(context.Background(), "old")
	if err != nil || p.Name != "Existing" || p.Status != project.InProgress || p.StartDate != "" || p.EndDate != "" {
		t.Fatalf("%+v %v", p, err)
	}
	h := lifecycleHandler(db, "2026-10-20T15:00:00Z")
	lifeRequest(t, h, "GET", "/api/projects/old/sprints/1", "", 200)
	legacy := lifeRequest(t, h, "POST", "/api/projects/legacy/finalize", `{"confirmed":true}`, 200)
	if legacy["end_date"] != nil {
		t.Fatal("invented historical date")
	}
	var count int
	if err = db.QueryRow(`SELECT count(*) FROM members WHERE id='member' AND project_id='old'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("lost member: %d %v", count, err)
	}
}
