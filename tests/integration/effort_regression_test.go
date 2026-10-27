package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite/migrations"
)

const estimateBody = `{"title":"Changed","description":"A","priority":"Alta","acceptance_criteria":["ok"],"estimated_hours":4.25}`

func execEffortSQL(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatal(err)
	}
}
func activeTotal(t *testing.T, db *sql.DB, story string) int {
	t.Helper()
	var v int
	if err := db.QueryRow(`SELECT COALESCE(SUM(hours_hundredths),0) FROM effort_entries WHERE backlog_item_id=? AND deleted_at IS NULL`, story).Scan(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

// BDD 01/02/03/04/06: explicitly reject missing estimates without corrupting legacy or points.
func TestUS16EstimatesAndRegularization(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	for _, value := range []string{"null", "0", "-1", "1.001", "1.500", "1e1", "\"1\""} {
		body := strings.Replace(estimateBody, "4.25", value, 1)
		expectCode(t, effortRequest(h, "POST", "/api/projects/p/backlog", body), 400)
		expectCode(t, effortRequest(h, "PUT", "/api/projects/p/backlog/a", body), 400)
	}
	missing := strings.Replace(estimateBody, `,"estimated_hours":4.25`, "", 1)
	expectCode(t, effortRequest(h, "POST", "/api/projects/p/backlog", missing), 400)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/backlog/a", missing), 400)
	for _, v := range []string{"0.01", "1", "1.50", "40.25"} {
		expectCode(t, effortRequest(h, "POST", "/api/projects/p/backlog", strings.Replace(estimateBody, "4.25", v, 1)), 201)
	}
	execEffortSQL(t, db, "UPDATE backlog_items SET story_points=5 WHERE id='a'")
	expectCode(t, effortRequest(h, "POST", effortPath, effortBody("5")), 201)
	w := effortRequest(h, "PUT", "/api/projects/p/backlog/a", estimateBody)
	expectCode(t, w, 200)
	v := readObject(t, w)
	if v["actual_hours"] != 5.0 || v["estimated_hours"] != 4.25 || v["story_points"] != 5.0 {
		t.Fatal(v)
	}
}

// BDD 17/18/40/43/44: finalized/completed edits, missing end, independent estimates.
func TestUS16FinalizedEffort(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	execEffortSQL(t, db, "UPDATE backlog_items SET status='Completada' WHERE id='a'")
	execEffortSQL(t, db, "UPDATE projects SET status='Finalizado',end_date='2026-10-20' WHERE id='p'")
	body := strings.Replace(effortBody("1"), "2026-10-25", "2026-10-20", 1)
	w := effortRequest(h, "POST", effortPath, body)
	expectCode(t, w, 201)
	id := readObject(t, w)["id"].(string)
	for _, date := range []string{"2026-10-21", "2026-09-30", "2026-10-27"} {
		expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(strings.Replace(body, "2026-10-20", date, 1), "}")+`,"backlog_item_id":"a"}`), 409)
	}
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(body, "}")+`,"backlog_item_id":"a"}`), 200)
	execEffortSQL(t, db, "UPDATE projects SET end_date=NULL WHERE id='p'")
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(body, "}")+`,"backlog_item_id":"a"}`), 409)
	expectCode(t, effortRequest(h, "GET", "/api/projects/p/effort/"+id, ""), 200)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/backlog/a", estimateBody), 200)
	expectCode(t, effortRequest(h, "POST", "/api/projects/p/effort/"+id+"/delete", `{"confirmed":true}`), 204)
	v := readObject(t, effortRequest(h, "GET", "/api/projects/p/backlog/a", ""))
	if v["status"] != "Completada" || v["actual_hours"] != 0.0 {
		t.Fatal(v)
	}
}

// BDD 20/21/22/23/26/46: updates subtract self, all fields, explicit Sprint correction.
func TestUS16ReplacementAndSprintValidation(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	execEffortSQL(t, db, `INSERT INTO sprints(id,project_id,goal)VALUES(1,'p','one'),(2,'p','two');
 INSERT INTO sprint_backlog_items(sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria)SELECT 1,'a',title,description,priority,status,acceptance_criteria FROM backlog_items WHERE id='a';
 INSERT INTO sprint_backlog_items(sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria)SELECT 2,'b',title,description,priority,status,acceptance_criteria FROM backlog_items WHERE id='b'`)
	w := effortRequest(h, "POST", effortPath, strings.TrimSuffix(effortBody("10"), "}")+`,"sprint_id":1}`)
	expectCode(t, w, 201)
	id := readObject(t, w)["id"].(string)
	expectCode(t, effortRequest(h, "POST", effortPath, effortBody("13.50")), 201)
	edit := strings.TrimSuffix(effortBody("10.50"), "}") + `,"backlog_item_id":"a"}`
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, edit), 200)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.Replace(edit, "10.50", "10.51", 1)), 409)
	move := strings.Replace(edit, `"backlog_item_id":"a"`, `"backlog_item_id":"b"`, 1)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, move), 409)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(move, "}")+`,"sprint_id":2}`), 200)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(move, "}")+`,"sprint_id":null}`), 200)
	if activeTotal(t, db, "a") != 1350 || activeTotal(t, db, "b") != 1050 {
		t.Fatal("incorrect totals after reassignment")
	}
	// Member/day/project are separate daily budgets.
	expectCode(t, effortRequest(h, "POST", effortPath, strings.Replace(effortBody("24"), `"m"`, `"n"`, 1)), 201)
	expectCode(t, effortRequest(h, "POST", effortPath, strings.Replace(effortBody("24"), "2026-10-25", "2026-10-24", 1)), 201)
	expectCode(t, effortRequest(h, "POST", "/api/projects/q/backlog/c/effort", strings.Replace(effortBody("24"), `"m"`, `"v"`, 1)), 201)
	edit = strings.Replace(edit, `"m"`, `"n"`, 1)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, edit), 409)
}

// BDD 31/36: trigger faults verify complete rollback, including propagated snapshots.
func TestUS16Rollback(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	w := effortRequest(h, "POST", effortPath, effortBody("2"))
	expectCode(t, w, 201)
	id := readObject(t, w)["id"].(string)
	for _, event := range []string{"INSERT", "UPDATE"} {
		execEffortSQL(t, db, `CREATE TRIGGER reject_effort BEFORE `+event+` ON effort_entries BEGIN SELECT RAISE(ABORT,'injected storage fault'); END`)
		if event == "INSERT" {
			expectCode(t, effortRequest(h, "POST", effortPath, effortBody("1")), 500)
		} else {
			expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(effortBody("3"), "}")+`,"backlog_item_id":"b"}`), 500)
			expectCode(t, effortRequest(h, "POST", "/api/projects/p/effort/"+id+"/delete", `{"confirmed":true}`), 500)
		}
		execEffortSQL(t, db, "DROP TRIGGER reject_effort")
		if activeTotal(t, db, "a") != 200 || activeTotal(t, db, "b") != 0 {
			t.Fatal("partial effort write")
		}
	}
	execEffortSQL(t, db, `INSERT INTO sprints(id,project_id,goal,closed)VALUES(1,'p','closed',1);
 INSERT INTO sprint_backlog_items(sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria)SELECT 1,id,title,description,priority,status,acceptance_criteria FROM backlog_items WHERE id='a';
 CREATE TRIGGER reject_snapshot BEFORE UPDATE ON sprint_backlog_items BEGIN SELECT RAISE(ABORT,'injected'); END`)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/backlog/a", estimateBody), 500)
	v := readObject(t, effortRequest(h, "GET", "/api/projects/p/backlog/a", ""))
	if v["title"] != "A" || v["estimated_hours"] != nil {
		t.Fatal(v)
	}
}

// BDD 30: run the actual previous migrations, preserve NULLs/identities, then reopen.
func TestUS16HistoricalMigrationAndReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	db, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	execEffortSQL(t, db, `CREATE TABLE carva_schema_migrations(version INTEGER PRIMARY KEY,name TEXT NOT NULL,applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)`)
	files, err := migrations.Files.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".sql") || name >= "000012" {
			continue
		}
		data, err := migrations.Files.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		execEffortSQL(t, db, string(data))
		var n int
		fmt.Sscanf(name, "%d_", &n)
		execEffortSQL(t, db, `INSERT INTO carva_schema_migrations(version,name)VALUES(?,?)`, n, name)
	}
	seedEffort(t, db)
	execEffortSQL(t, db, `INSERT INTO sprints(id,project_id,goal,closed)VALUES(1,'p','archived',1),(2,'p','open',0);
 INSERT INTO sprint_backlog_items(sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria)
 SELECT s.id,'a','Original','Original','Alta','Pendiente','["original"]' FROM sprints s WHERE s.project_id='p'`)
	for range 2 {
		if err := sqlite.Migrate(context.Background(), db); err != nil {
			t.Fatal(err)
		}
	}
	h := effortHandler(db)
	w := effortRequest(h, "POST", effortPath, effortBody("1.25"))
	expectCode(t, w, 201)
	id := readObject(t, w)["id"].(string)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/backlog/a", estimateBody), 200)
	expectCode(t, effortRequest(h, "POST", "/api/projects/p/effort/"+id+"/delete", `{"confirmed":true}`), 204)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	v := readObject(t, effortRequest(effortHandler(db), "GET", "/api/projects/p/backlog/a", ""))
	if v["estimated_hours"] != 4.25 || v["actual_hours"] != 0.0 {
		t.Fatal(v)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM effort_entries WHERE id=? AND deleted_at IS NOT NULL`, id).Scan(&count); err != nil || count != 1 {
		t.Fatal(count, err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM sprint_backlog_items WHERE backlog_item_id='a' AND estimated_hours_hundredths=425 AND title='Original' AND acceptance_criteria='["original"]'`).Scan(&count); err != nil || count != 2 {
		t.Fatal(count, err)
	}
}

func concurrentEffort(t *testing.T, ops ...func() int) []int {
	t.Helper()
	start := make(chan struct{})
	var wg sync.WaitGroup
	codes := make([]int, len(ops))
	for i, op := range ops {
		wg.Add(1)
		go func(i int, op func() int) { defer wg.Done(); <-start; codes[i] = op() }(i, op)
	}
	close(start)
	wg.Wait()
	return codes
}

// BDD 33/35/41/42/45: competing independent pools never leave inconsistent rows.
func TestUS16ConcurrentLifecycle(t *testing.T) {
	for _, operation := range []string{"start", "delete-start", "delete-edit", "finalize", "edit-budget"} {
		t.Run(operation, func(t *testing.T) {
			path, db := openMigratedProjectDB(t)
			seedEffort(t, db)
			other, err := sqlite.OpenPath(context.Background(), path)
			if err != nil {
				t.Fatal(err)
			}
			defer other.Close()
			h, h2 := effortHandler(db), effortHandler(other)
			var codes []int
			switch operation {
			case "start":
				body := strings.Replace(effortBody("1"), "2026-10-25", "2026-10-05", 1)
				codes = concurrentEffort(t, func() int { return effortRequest(h, "POST", effortPath, body).Code }, func() int { return effortRequest(h2, "PATCH", "/api/projects/p", `{"start_date":"2026-10-10"}`).Code })
				if (codes[0] == 201) == (codes[1] == 200) {
					t.Fatal(codes)
				}
			case "delete-start", "delete-edit":
				w := effortRequest(h, "POST", effortPath, strings.Replace(effortBody("1"), "2026-10-25", "2026-10-05", 1))
				expectCode(t, w, 201)
				id := readObject(t, w)["id"].(string)
				codes = concurrentEffort(t, func() int {
					return effortRequest(h, "POST", "/api/projects/p/effort/"+id+"/delete", `{"confirmed":true}`).Code
				}, func() int {
					if operation == "delete-start" {
						return effortRequest(h2, "PATCH", "/api/projects/p", `{"start_date":"2026-10-10"}`).Code
					}
					return effortRequest(h2, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(effortBody("3"), "}")+`,"backlog_item_id":"a"}`).Code
				})
				if codes[0] != 204 || (codes[1] != 200 && codes[1] != 409 && codes[1] != 404) {
					t.Fatal(codes)
				}
				if activeTotal(t, db, "a") != 0 {
					t.Fatal("resurrected")
				}
			case "finalize":
				codes = concurrentEffort(t, func() int { return effortRequest(h, "POST", effortPath, effortBody("1")).Code }, func() int { return effortRequest(h2, "POST", "/api/projects/p/finalize", `{"confirmed":true}`).Code })
				if codes[0] != 201 || codes[1] != 200 {
					t.Fatal(codes)
				}
			case "edit-budget":
				expectCode(t, effortRequest(h, "POST", effortPath, effortBody("20")), 201)
				ids := []string{}
				for _, day := range []string{"2026-10-23", "2026-10-24"} {
					w := effortRequest(h, "POST", effortPath, strings.Replace(effortBody("1"), "2026-10-25", day, 1))
					expectCode(t, w, 201)
					ids = append(ids, readObject(t, w)["id"].(string))
				}
				body := strings.TrimSuffix(effortBody("3"), "}") + `,"backlog_item_id":"a"}`
				codes = concurrentEffort(t, func() int { return effortRequest(h, "PUT", "/api/projects/p/effort/"+ids[0], body).Code }, func() int { return effortRequest(h2, "PUT", "/api/projects/p/effort/"+ids[1], body).Code })
				if !((codes[0] == 200 && codes[1] == 409) || (codes[1] == 200 && codes[0] == 409)) {
					t.Fatal(codes)
				}
			}
			var invalid int
			if err := db.QueryRow(`SELECT COUNT(*) FROM effort_entries e JOIN projects p ON p.id=e.project_id WHERE e.deleted_at IS NULL AND (e.work_date<p.start_date OR e.work_date>COALESCE(p.end_date,'9999-12-31'))`).Scan(&invalid); err != nil || invalid != 0 {
				t.Fatal(invalid, err)
			}
		})
	}
}

// BDD 29/34/37: use actual application operations for assignment, transfer, agreement and completion.
func TestUS16ConcurrentPlanning(t *testing.T) {
	for _, operation := range []string{"assign", "transfer", "transfer-existing", "agreement", "complete"} {
		t.Run(operation, func(t *testing.T) {
			path, db := openMigratedProjectDB(t)
			seedEffort(t, db)
			other, err := sqlite.OpenPath(context.Background(), path)
			if err != nil {
				t.Fatal(err)
			}
			defer other.Close()
			execEffortSQL(t, db, `INSERT INTO sprints(id,project_id,goal)VALUES(1,'p','origin'),(2,'p','destination')`)
			br := sqlite.NewBacklogRepository(db)
			ss := sprints.NewService(sqlite.NewSprintRepository(db), br)
			if operation != "assign" {
				if _, err := ss.AssignBacklogItem(context.Background(), "p", 1, "a"); err != nil {
					t.Fatal(err)
				}
			}
			if operation != "assign" {
				w := effortRequest(effortHandler(db), "POST", effortPath, strings.TrimSuffix(effortBody("1"), "}")+`,"sprint_id":1}`)
				expectCode(t, w, 201)
			}
			if operation == "transfer-existing" {
				if _, err := ss.AssignBacklogItem(context.Background(), "p", 2, "a"); err != nil {
					t.Fatal(err)
				}
			}
			ss = sprints.NewService(sqlite.NewSprintRepository(other), sqlite.NewBacklogRepository(other))
			execEffortSQL(t, db, `INSERT INTO estimation_rounds VALUES('r','p','a',1,'finalized'); INSERT INTO estimation_round_participants VALUES('r','m'); INSERT INTO estimation_votes VALUES('r','m',5)`)
			codes := concurrentEffort(t, func() int {
				return effortRequest(effortHandler(db), "PUT", "/api/projects/p/backlog/a", estimateBody).Code
			}, func() int {
				var err error
				switch operation {
				case "assign":
					_, err = ss.AssignBacklogItem(context.Background(), "p", 1, "a")
				case "transfer", "transfer-existing":
					_, err = ss.Close(context.Background(), "p", 1)
				case "complete":
					_, err = ss.CompleteBacklogItem(context.Background(), "p", 1, "a", true)
				case "agreement":
					_, err = sqlite.NewPlanningPokerRepository(other).RegisterAgreement(context.Background(), "p", "a", "r", "m")
				}
				if err != nil {
					return 500
				}
				return 200
			})
			if codes[0] != 200 || codes[1] != 200 {
				t.Fatal(codes)
			}
			var bad int
			if err := db.QueryRow(`SELECT COUNT(*) FROM sprint_backlog_items WHERE backlog_item_id='a' AND (estimated_hours_hundredths IS NULL OR estimated_hours_hundredths<>425)`).Scan(&bad); err != nil || bad != 0 {
				t.Fatal(bad, err)
			}
			if operation != "assign" {
				var associated int
				if err := db.QueryRow(`SELECT sprint_id FROM effort_entries WHERE project_id='p' AND backlog_item_id='a'`).Scan(&associated); err != nil || associated != 1 || activeTotal(t, db, "a") != 100 {
					t.Fatal(associated, err)
				}
			}
			value, err := appbacklog.NewService(br).Get(context.Background(), "p", "a")
			if err != nil {
				t.Fatal(err)
			}
			if operation == "complete" && value.Status != "Completada" {
				t.Fatal(value)
			}
			if operation == "agreement" && (value.StoryPoints == nil || *value.StoryPoints != 5) {
				t.Fatal(value)
			}
		})
	}
}

// BDD 38: the public endpoint reports bounded write contention without changes.
func TestUS16BoundedContention(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(context.Background(), "BEGIN IMMEDIATE"); err != nil {
		t.Fatal(err)
	}
	defer conn.ExecContext(context.Background(), "ROLLBACK")
	start := time.Now()
	expectCode(t, effortRequest(effortHandler(db), "POST", effortPath, effortBody("1")), 503)
	if time.Since(start) > 8*time.Second {
		t.Fatal("unbounded wait")
	}
}

// BDD 10/12/27: every scope is resolved by identity in the route's project.
func TestUS16ScopedResources(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	execEffortSQL(t, db, `INSERT INTO sprints(id,project_id,goal)VALUES(1,'p','unassigned'),(2,'q','foreign')`)
	for _, tc := range []struct {
		path, body string
		want       int
	}{
		{effortPath, strings.Replace(effortBody("1"), `"m"`, `"missing"`, 1), 404},
		{effortPath, strings.TrimSuffix(effortBody("1"), "}") + `,"sprint_id":1}`, 409},
		{effortPath, strings.TrimSuffix(effortBody("1"), "}") + `,"sprint_id":2}`, 404},
		{effortPath, strings.TrimSuffix(effortBody("1"), "}") + `,"sprint_id":999}`, 404},
		{"/api/projects/p/backlog/c/effort", effortBody("1"), 404},
		{"/api/projects/missing/backlog/a/effort", effortBody("1"), 404},
		{effortPath, strings.TrimSuffix(effortBody("1"), "}") + `,"project_id":"q"}`, 400},
	} {
		expectCode(t, effortRequest(h, "POST", tc.path, tc.body), tc.want)
	}
	w := effortRequest(h, "POST", effortPath, effortBody("1"))
	expectCode(t, w, 201)
	id := readObject(t, w)["id"].(string)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, strings.TrimSuffix(effortBody("1"), "}")+`,"backlog_item_id":"c"}`), 404)
	expectCode(t, effortRequest(h, "POST", "/api/projects/q/effort/"+id+"/delete", `{"confirmed":true}`), 404)
	if activeTotal(t, db, "a") != 100 {
		t.Fatal("cross-project operation wrote data")
	}
}
