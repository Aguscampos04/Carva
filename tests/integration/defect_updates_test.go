package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Aguscampos04/Carva/internal/application/projects"
	domain "github.com/Aguscampos04/Carva/internal/domain/defect"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

type defectUpdateFixture struct {
	db      *sql.DB
	path    string
	handler http.Handler
	url, id string
}

func newDefectUpdateFixture(t *testing.T) defectUpdateFixture {
	t.Helper()
	path, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	if _, err := db.Exec(`INSERT INTO sprints(id,project_id,goal,closed) VALUES(1,'p','Detection',1),(2,'p','Resolution',0),(3,'q','Foreign',0)`); err != nil {
		t.Fatal(err)
	}
	h := defectHandler(db)
	created := defectRequest(h, "POST", "/api/projects/p/backlog/a/defects", `{"description":"Original bug","severity":"Alta","status":"Abierto"}`)
	expectCode(t, created, 201)
	id := readObject(t, created)["id"].(string)
	return defectUpdateFixture{db: db, path: path, handler: h, id: id, url: "/api/projects/p/backlog/a/defects/" + id}
}
func listedDefect(t *testing.T, h http.Handler) map[string]any {
	t.Helper()
	r := defectRequest(h, "GET", "/api/projects/p/backlog/a/defects", "")
	expectCode(t, r, 200)
	var items []map[string]any
	if err := json.Unmarshal(r.Body.Bytes(), &items); err != nil || len(items) != 1 {
		t.Fatalf("%s %v", r.Body, err)
	}
	return items[0]
}

// AC-01/02/03, BDD-01/02/03/04/08/09/10: nullable partial updates persist.
func TestUS19UpdatesPersistAndPreserveOmittedFields(t *testing.T) {
	f := newDefectUpdateFixture(t)
	expectCode(t, defectRequest(f.handler, "PATCH", f.url, `{"detection_sprint_id":1}`), 200)
	expectCode(t, defectRequest(f.handler, "PATCH", f.url, `{"status":"Resuelto","resolution_sprint_id":2}`), 200)
	for _, status := range []string{"Cerrado", "Abierto", "En progreso", "Resuelto"} {
		r := defectRequest(f.handler, "PATCH", f.url, `{"status":"`+status+`"}`)
		expectCode(t, r, 200)
		got := readObject(t, r)
		if got["id"] != f.id || got["description"] != "Original bug" || got["severity"] != "Alta" || got["detection_sprint_id"] != float64(1) || got["resolution_sprint_id"] != float64(2) || got["status"] != status {
			t.Fatal(got)
		}
	}
	if err := f.db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err := sqlite.OpenPath(context.Background(), f.path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f.handler = defectHandler(db)
	got := listedDefect(t, f.handler)
	if got["status"] != "Resuelto" || got["resolution_sprint_id"] != float64(2) {
		t.Fatal(got)
	}
	expectCode(t, defectRequest(f.handler, "PATCH", f.url, `{"resolution_sprint_id":null}`), 200)
	got = listedDefect(t, f.handler)
	if got["resolution_sprint_id"] != nil || got["status"] != "Resuelto" || got["detection_sprint_id"] != float64(1) {
		t.Fatal(got)
	}
	if _, err := db.Exec(`UPDATE projects SET status='Finalizado',end_date='2026-10-27' WHERE id='p'; UPDATE backlog_items SET status='Completada' WHERE id='a'`); err != nil {
		t.Fatal(err)
	}
	expectCode(t, defectRequest(f.handler, "PATCH", f.url, `{"resolution_sprint_id":1,"detection_sprint_id":null}`), 200)
	got = listedDefect(t, f.handler)
	if got["detection_sprint_id"] != nil || got["resolution_sprint_id"] != float64(1) || got["status"] != "Resuelto" {
		t.Fatal(got)
	}
	// Repeat the same change: no additional record, identity unchanged.
	expectCode(t, defectRequest(f.handler, "PATCH", f.url, `{"resolution_sprint_id":1}`), 200)
	if listedDefect(t, f.handler)["id"] != f.id {
		t.Fatal("identity changed")
	}
}

// BDD-05/06/07/11: reject malformed/foreign writes atomically, with safe errors.
func TestUS19RejectsInvalidAndForeignUpdates(t *testing.T) {
	f := newDefectUpdateFixture(t)
	before := listedDefect(t, f.handler)
	for _, body := range []string{`{}`, `null`, `[]`, `{"status":null}`, `{"status":"Unknown"}`, `{"status":""}`, `{"status":true}`, `{"resolution_sprint_id":0}`, `{"detection_sprint_id":-1}`, `{"detection_sprint_id":1.5}`, `{"detection_sprint_id":"1"}`, `{"status":"Abierto","status":"Resuelto"}`, `{"description":"Other"}`, `{"id":"other"}`, `{"status":"Resuelto"} {}`, `{"Status":"Resuelto"}`} {
		expectCode(t, defectRequest(f.handler, "PATCH", f.url, body), 400)
		if !reflect.DeepEqual(before, listedDefect(t, f.handler)) {
			t.Fatal("invalid input changed data", body)
		}
	}
	for _, body := range []string{`{"status":"Resuelto","resolution_sprint_id":3}`, `{"detection_sprint_id":999999}`} {
		expectCode(t, defectRequest(f.handler, "PATCH", f.url, body), 404)
		if !reflect.DeepEqual(before, listedDefect(t, f.handler)) {
			t.Fatal("partial update", body)
		}
	}
	for _, url := range []string{strings.Replace(f.url, "/a/", "/b/", 1), strings.Replace(f.url, "/p/", "/q/", 1), f.url + "missing", strings.Replace(f.url, "/p/", "/missing/", 1)} {
		expectCode(t, defectRequest(f.handler, "PATCH", url, `{"status":"Resuelto"}`), 404)
	}
	if _, err := f.db.Exec(`CREATE TRIGGER reject_defect_update BEFORE UPDATE ON defects BEGIN SELECT RAISE(ABORT,'controlled secret'); END`); err != nil {
		t.Fatal(err)
	}
	r := defectRequest(f.handler, "PATCH", f.url, `{"status":"Resuelto","detection_sprint_id":1}`)
	expectCode(t, r, 500)
	if strings.Contains(r.Body.String(), "controlled secret") || !reflect.DeepEqual(before, listedDefect(t, f.handler)) {
		t.Fatal("leaked error or partial data")
	}
	if err := f.db.Close(); err != nil {
		t.Fatal(err)
	}
	expectCode(t, defectRequest(f.handler, "PATCH", f.url, `{"status":"Resuelto"}`), 500)
}

// BDD-12: two pools update independent fields against the latest locked row.
func TestUS19ConcurrentUpdatesPreserveDistinctFields(t *testing.T) {
	f := newDefectUpdateFixture(t)
	other, err := sqlite.OpenPath(context.Background(), f.path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	for round := 0; round < 8; round++ {
		expectCode(t, defectRequest(f.handler, "PATCH", f.url, `{"status":"Abierto","detection_sprint_id":null}`), 200)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for _, op := range []struct {
			h    http.Handler
			body string
		}{{f.handler, `{"status":"Resuelto"}`}, {defectHandler(other), `{"detection_sprint_id":1}`}} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				r := defectRequest(op.h, "PATCH", f.url, op.body)
				if r.Code != 200 {
					t.Errorf("%d %s", r.Code, r.Body)
				}
			}()
		}
		close(start)
		wg.Wait()
		got := listedDefect(t, f.handler)
		if got["status"] != "Resuelto" || got["detection_sprint_id"] != float64(1) {
			t.Fatal("lost update", got)
		}
	}
}

func TestUS19UpdateHonorsContentionDeadline(t *testing.T) {
	f := newDefectUpdateFixture(t)
	conn, err := f.db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(context.Background(), "BEGIN IMMEDIATE"); err != nil {
		t.Fatal(err)
	}
	defer conn.ExecContext(context.Background(), "ROLLBACK")
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	status := domain.StatusResolved
	_, err = sqlite.NewDefectRepository(f.db).Update(ctx, "p", "a", f.id, domain.Patch{Status: &status})
	if !errors.Is(err, projects.ErrBusy) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("contention error=%v", err)
	}
}
