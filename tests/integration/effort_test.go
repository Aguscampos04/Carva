package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	appeffort "github.com/Aguscampos04/Carva/internal/application/effort"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func effortHandler(db *sql.DB) http.Handler {
	clock := func() time.Time { return time.Date(2026, 10, 26, 15, 0, 0, 0, time.UTC) }
	br := sqlite.NewBacklogRepository(db)
	base := httpapi.NewHandler(projects.NewServiceWithClock(sqlite.NewProjectRepository(db), clock), appbacklog.NewService(br), sprints.NewService(sqlite.NewSprintRepository(db), br))
	return httpapi.WithEffort(base, appeffort.NewService(sqlite.NewEffortRepository(db), clock))
}
func seedEffort(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, q := range []string{
		`INSERT INTO projects(id,name,status,start_date) VALUES('p','Carva Demo','En progreso','2026-10-01'),('q','Vecino','Planificado',NULL)`,
		`INSERT INTO members(id,project_id,name) VALUES('m','p','Ana'),('n','p','Ana'),('v','q','Ana')`,
		`INSERT INTO backlog_items(id,project_id,title,description,priority,status,acceptance_criteria) VALUES('a','p','A','A','Alta','Pendiente','["ok"]'),('b','p','B','B','Alta','Pendiente','["ok"]'),('c','q','C','C','Alta','Pendiente','["ok"]')`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
}
func effortRequest(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}
func effortBody(hours string) string {
	return `{"member_id":"m","work_date":"2026-10-25","activity":"Review","hours":` + hours + `}`
}
func expectCode(t *testing.T, w *httptest.ResponseRecorder, code int) {
	t.Helper()
	if w.Code != code {
		t.Fatalf("status %d want %d: %s", w.Code, code, w.Body)
	}
}
func readObject(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var v map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	return v
}

// AC-US17-01/02: comparison keeps estimates, active effort and Story Points distinct.
func TestUS17EffortComparisonUsesOnlyRelatedActiveRecords(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	if _, err := db.Exec(`UPDATE backlog_items SET estimated_hours_hundredths=250,story_points=5 WHERE project_id='p' AND id='a'`); err != nil {
		t.Fatal(err)
	}
	h := effortHandler(db)

	for _, tc := range []struct {
		story string
		hours string
	}{
		{story: "a", hours: "1.25"},
		{story: "a", hours: "2.50"},
		{story: "a", hours: "4.00"},
		{story: "b", hours: "7.00"},
	} {
		response := effortRequest(h, "POST", "/api/projects/p/backlog/"+tc.story+"/effort", effortBody(tc.hours))
		expectCode(t, response, http.StatusCreated)
	}

	var removable string
	if err := db.QueryRow(`SELECT id FROM effort_entries WHERE backlog_item_id='a' AND hours_hundredths=400`).Scan(&removable); err != nil {
		t.Fatal(err)
	}
	expectCode(t, effortRequest(h, "POST", "/api/projects/p/effort/"+removable+"/delete", `{"confirmed":true}`), http.StatusNoContent)

	storyA := readObject(t, effortRequest(h, "GET", "/api/projects/p/backlog/a", ""))
	if storyA["estimated_hours"] != 2.5 || storyA["actual_hours"] != 3.75 || storyA["story_points"] != float64(5) {
		t.Fatalf("comparison must retain separate hour sources and Story Points: %#v", storyA)
	}
	storyB := readObject(t, effortRequest(h, "GET", "/api/projects/p/backlog/b", ""))
	if storyB["actual_hours"] != float64(7) {
		t.Fatalf("effort from another story leaked into comparison: %#v", storyB)
	}
}

const effortPath = "/api/projects/p/backlog/a/effort"

// BDD 05/07/11/19/22/24/27/30/44: legacy, independent entries, totals and soft delete.
func TestUS16EffortLifecycle(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	w := effortRequest(h, "POST", effortPath, effortBody("1.25"))
	expectCode(t, w, 201)
	id := readObject(t, w)["id"].(string)
	w = effortRequest(h, "POST", effortPath, effortBody("2.50"))
	expectCode(t, w, 201)
	if readObject(t, w)["id"] == id {
		t.Fatal("deduplicated entries")
	}
	detail := readObject(t, effortRequest(h, "GET", "/api/projects/p/backlog/a", ""))
	if detail["estimated_hours"] != nil || detail["actual_hours"] != 3.75 {
		t.Fatal(detail)
	}
	edit := `{"backlog_item_id":"b","member_id":"n","work_date":"2026-10-24","activity":"Corrected","hours":4}`
	w = effortRequest(h, "PUT", "/api/projects/p/effort/"+id, edit)
	expectCode(t, w, 200)
	if readObject(t, w)["id"] != id {
		t.Fatal("changed identity")
	}
	expectCode(t, effortRequest(h, "GET", "/api/projects/q/effort/"+id, ""), 404)
	expectCode(t, effortRequest(h, "POST", "/api/projects/p/effort/"+id+"/delete", `{"confirmed":false}`), 400)
	for i := 0; i < 2; i++ {
		expectCode(t, effortRequest(h, "POST", "/api/projects/p/effort/"+id+"/delete", `{"confirmed":true}`), 204)
	}
	expectCode(t, effortRequest(h, "GET", "/api/projects/p/effort/"+id, ""), 404)
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/effort/"+id, edit), 404)
	var date string
	if err := db.QueryRow("SELECT work_date FROM effort_entries WHERE id=? AND deleted_at IS NOT NULL", id).Scan(&date); err != nil || date != "2026-10-24" {
		t.Fatal(date, err)
	}
}

// BDD 08/09/10/12/14/18/20/43.
func TestUS16InvalidEffortAndDates(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	for _, tc := range []struct {
		body string
		code int
	}{
		{effortBody("0"), 400}, {effortBody("-1"), 400}, {effortBody("1.001"), 400}, {effortBody("1e1"), 400}, {effortBody("24.01"), 409},
		{strings.Replace(effortBody("1"), `"m"`, `"v"`, 1), 404},
		{strings.Replace(effortBody("1"), "2026-10-25", "2026-09-30", 1), 409},
		{strings.Replace(effortBody("1"), "2026-10-25", "2026-10-27", 1), 409},
		{strings.Replace(effortBody("1"), "2026-10-25", "2026-02-30", 1), 400},
		{strings.Replace(effortBody("1"), "Review", "   ", 1), 400},
		{`{}`, 400}, {`null`, 400}, {`[]`, 400}, {effortBody("1") + ` {}`, 400},
		{strings.Replace(effortBody("1"), `"hours":1`, `"hours":1,"hours":2`, 1), 400},
	} {
		w := effortRequest(h, "POST", effortPath, tc.body)
		if w.Code != tc.code {
			t.Errorf("%s: %d want %d %s", tc.body, w.Code, tc.code, w.Body)
		}
	}
	for _, status := range []string{"Planificado", "En progreso", "Finalizado"} {
		if _, err := db.Exec("UPDATE projects SET status=?,end_date=NULL WHERE id='p'", status); err != nil {
			t.Fatal(err)
		}
		want := 201
		if status == "Finalizado" {
			want = 409
		}
		expectCode(t, effortRequest(h, "POST", effortPath, effortBody("0.01")), want)
	}
}

// BDD 20/21/23/32/42: independent connections, atomic daily limit.
func TestUS16ConcurrentDailyLimit(t *testing.T) {
	path, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	expectCode(t, effortRequest(h, "POST", effortPath, effortBody("23")), 201)
	other, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	handlers := []http.Handler{h, effortHandler(other)}
	codes := make(chan int, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go func(h http.Handler) {
			defer wg.Done()
			<-start
			codes <- effortRequest(h, "POST", effortPath, effortBody("1")).Code
		}(handler)
	}
	close(start)
	wg.Wait()
	close(codes)
	counts := map[int]int{}
	for c := range codes {
		counts[c]++
	}
	if counts[201] != 1 || counts[409] != 1 {
		t.Fatal(counts)
	}
	var total int
	if err := db.QueryRow("SELECT SUM(hours_hundredths) FROM effort_entries WHERE deleted_at IS NULL").Scan(&total); err != nil || total != 2400 {
		t.Fatal(total, err)
	}
}

// BDD 06/15/16/26/28/29/36/45/46.
func TestUS16SnapshotsSprintAndStart(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	h := effortHandler(db)
	for _, q := range []string{
		`INSERT INTO sprints(id,project_id,goal,closed) VALUES(1,'p','old',1),(2,'p','new',0)`,
		`INSERT INTO sprint_backlog_items(sprint_id,backlog_item_id,title,description,priority,status,acceptance_criteria) SELECT 1,id,title,description,priority,status,acceptance_criteria FROM backlog_items WHERE id IN ('a','b')`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	body := `{"title":"Changed","description":"A","priority":"Alta","acceptance_criteria":["ok"],"estimated_hours":4.25}`
	expectCode(t, effortRequest(h, "PUT", "/api/projects/p/backlog/a", body), 200)
	var estimate int
	var title string
	if err := db.QueryRow("SELECT estimated_hours_hundredths,title FROM sprint_backlog_items WHERE backlog_item_id='a'").Scan(&estimate, &title); err != nil || estimate != 425 || title != "A" {
		t.Fatal(estimate, title, err)
	}
	w := effortRequest(h, "POST", effortPath, strings.TrimSuffix(effortBody("1"), "}")+`,"sprint_id":1}`)
	expectCode(t, w, 201)
	id := readObject(t, w)["id"].(string)
	edit := `{"backlog_item_id":"b","member_id":"m","work_date":"2026-10-25","activity":"move","hours":1}`
	w = effortRequest(h, "PUT", "/api/projects/p/effort/"+id, edit)
	expectCode(t, w, 200)
	if readObject(t, w)["sprint_id"] != float64(1) {
		t.Fatal("lost compatible Sprint")
	}
	expectCode(t, effortRequest(h, "PATCH", "/api/projects/p", `{"name":"bad","start_date":"2026-10-26"}`), 409)
	expectCode(t, effortRequest(h, "POST", "/api/projects/p/effort/"+id+"/delete", `{"confirmed":true}`), 204)
	expectCode(t, effortRequest(h, "PATCH", "/api/projects/p", `{"start_date":"2026-10-26"}`), 200)
	if _, err := db.Exec("UPDATE projects SET status='Finalizado',end_date=NULL WHERE id='p'"); err != nil {
		t.Fatal(err)
	}
	expectCode(t, effortRequest(h, "GET", "/api/projects/p/backlog/a/effort-context", ""), 200)
	expectCode(t, effortRequest(h, "POST", effortPath, effortBody("1")), 409)
}
