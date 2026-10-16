package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/members"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/domain/member"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
	"github.com/google/uuid"
)

func memberHandler(db *sql.DB) http.Handler {
	return httpapi.NewHandler(projects.NewService(sqlite.NewProjectRepository(db)), nil, nil,
		members.NewService(sqlite.NewMemberRepository(db)))
}

func memberRequest(h http.Handler, method, projectID, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, "/api/projects/"+projectID+"/members", strings.NewReader(body)))
	return w
}

func createMemberProject(t *testing.T, db *sql.DB, name string) string {
	t.Helper()
	p, err := projects.NewService(sqlite.NewProjectRepository(db)).Create(context.Background(), name)
	if err != nil {
		t.Fatal(err)
	}
	return p.ID
}

func decodeMember(t *testing.T, w *httptest.ResponseRecorder) member.Member {
	t.Helper()
	if w.Code != http.StatusCreated {
		t.Fatalf("POST = %d %s", w.Code, w.Body)
	}
	var result member.Member
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if _, err := uuid.Parse(result.ID); err != nil {
		t.Fatal(err)
	}
	return result
}

func listMembers(t *testing.T, h http.Handler, p string) []member.Member {
	t.Helper()
	w := memberRequest(h, "GET", p, "")
	if w.Code != http.StatusOK {
		t.Fatalf("GET = %d %s", w.Code, w.Body)
	}
	var result []member.Member
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("list must not be null")
	}
	return result
}

// BDD-US-03-01/04/05/06/09/14: actual REST + SQLite, no mock persistence.
func TestMembersHTTPRegistrationIsolationAndReopen(t *testing.T) {
	path, db := openMigratedProjectDB(t)
	p := createMemberProject(t, db, "Carva Demo")
	other := createMemberProject(t, db, "Carva Alternativo")
	h := memberHandler(db)
	if got := listMembers(t, h, p); len(got) != 0 {
		t.Fatal(got)
	}
	w := memberRequest(h, "POST", p, `{"name":"  Ana  Pérez  ","id":"chosen","project_id":"`+other+`","email":"unused"}`)
	a := decodeMember(t, w)
	var fields map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 3 || a.ID == "chosen" || a.ProjectID != p || a.Name != "Ana  Pérez" {
		t.Fatal(fields)
	}
	b := decodeMember(t, memberRequest(h, "POST", p, `{"name":"Ana  Pérez"}`))
	c := decodeMember(t, memberRequest(h, "POST", other, `{"name":"Ana  Pérez"}`))
	if a.ID == b.ID || a.ID == c.ID || b.ID == c.ID {
		t.Fatal("identities must differ")
	}
	if got := listMembers(t, h, p); len(got) != 2 {
		t.Fatal(got)
	}
	if got := listMembers(t, h, other); !reflect.DeepEqual(got, []member.Member{c}) {
		t.Fatal(got)
	}
	project, err := projects.NewService(sqlite.NewProjectRepository(db)).Get(context.Background(), p)
	if err != nil || string(project.Status) != "Planificado" {
		t.Fatalf("project changed: %v %v", project, err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := sqlite.OpenPath(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	got := listMembers(t, memberHandler(reopened), p)
	byID := map[string]member.Member{}
	for _, m := range got {
		byID[m.ID] = m
	}
	if len(got) != 2 || byID[a.ID] != a || byID[b.ID] != b {
		t.Fatalf("reopen changed identity: %v", got)
	}
}

// BDD-US-03-02/03/08/17/19: reject the entire invalid input, before project lookup.
func TestMembersHTTPStrictJSONAndNameValidation(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	p := createMemberProject(t, db, "Carva Demo")
	h := memberHandler(db)
	invalid := []string{"", `null`, `[]`, `"Ana"`, `42`, `true`, `{"name":`, `{} {}`, `{"name":"Ana"} trailing`,
		`{}`, `{"name":null}`, `{"name":""}`, `{"name":"  "}`, `{"name":"\t\n\u00a0"}`,
		`{"name":42}`, `{"name":true}`, `{"name":[]}`, `{"name":{}}`,
		`{"name":"Ana","name":"Bruno"}`, `{"Name":"Ana"}`, `{"extra":1,"extra":2,"name":"Ana"}`,
		"{\"name\":\"\xff\"}",
	}
	for _, body := range invalid {
		t.Run(body, func(t *testing.T) {
			for _, projectID := range []string{p, "missing"} {
				w := memberRequest(h, "POST", projectID, body)
				if w.Code != 400 {
					t.Fatalf("body %q = %d %s", body, w.Code, w.Body)
				}
				var failure map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &failure); err != nil || failure["error"] == "" {
					t.Fatalf("error: %v %s", err, w.Body)
				}
			}
		})
	}
	if got := listMembers(t, h, p); len(got) != 0 {
		t.Fatal("invalid writes persisted", got)
	}
	for _, tc := range []struct{ body, want string }{
		{`{"name":"\t\n\u00a0Ana  Pérez\u00a0\n\t"}`, "Ana  Pérez"},
		{`{"name":"a"}`, "a"}, {`{"name":"<b>Ana</b>"}`, "<b>Ana</b>"},
		{`{"name":"\u0000Ana"}`, "\x00Ana"},
	} {
		got := decodeMember(t, memberRequest(h, "POST", p, tc.body))
		if got.Name != tc.want {
			t.Fatal(got)
		}
	}
}

// BDD-US-03-10/11/12/13: missing project versus actual storage failures.
func TestMembersHTTPMissingAndStorageErrors(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	p := createMemberProject(t, db, "Carva Demo")
	h := memberHandler(db)
	for _, method := range []string{"POST", "GET"} {
		w := memberRequest(h, method, "00000000-0000-0000-0000-000000000999", `{"name":"Ana"}`)
		if w.Code != 404 {
			t.Fatalf("missing: %d %s", w.Code, w.Body)
		}
	}
	// Dropping only members keeps project lookup working and exercises save/list failures.
	if _, err := db.Exec(`DROP TABLE members`); err != nil {
		t.Fatal(err)
	}
	for _, method := range []string{"POST", "GET"} {
		w := memberRequest(h, method, p, `{"name":"Ana"}`)
		if w.Code != 500 || strings.Contains(w.Body.String(), "SQL") || strings.Contains(w.Body.String(), "table") {
			t.Fatalf("storage error: %d %s", w.Code, w.Body)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if w := memberRequest(h, "GET", p, ""); w.Code != 500 {
		t.Fatal(w.Code)
	}
}

func TestMembersHTTPDoesNotExposeEditOrDelete(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	p := createMemberProject(t, db, "Carva Demo")
	h := memberHandler(db)
	m := decodeMember(t, memberRequest(h, "POST", p, `{"name":"Ana"}`))
	for _, method := range []string{"PUT", "PATCH", "DELETE"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, "/api/projects/"+p+"/members/"+m.ID, strings.NewReader(`{"name":"Bruno"}`)))
		if w.Code < 400 {
			t.Fatal("unexpected endpoint", method, w.Code)
		}
	}
	if got := listMembers(t, h, p); !reflect.DeepEqual(got, []member.Member{m}) {
		t.Fatal(got)
	}
}

// BDD-US-03-18. A failed concurrent write may be reported as such; successful ones cannot overwrite.
func TestMembersConcurrentHomonyms(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	db.SetMaxOpenConns(1)
	p := createMemberProject(t, db, "Carva Demo")
	h := memberHandler(db)
	results := make(chan *httptest.ResponseRecorder, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); results <- memberRequest(h, "POST", p, `{"name":"Ana"}`) }()
	}
	wg.Wait()
	close(results)
	for w := range results {
		decodeMember(t, w)
	}
	got := listMembers(t, h, p)
	if len(got) != 2 || got[0].ID == got[1].ID {
		t.Fatal(got)
	}
}

func TestMembersSQLiteForeignKeysOnEveryConnection(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	ctx := context.Background()
	db.SetMaxOpenConns(3)
	var conns []*sql.Conn
	for range 3 {
		conn, err := db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		conns = append(conns, conn)
		defer conn.Close()
		var enabled int
		if err := conn.QueryRowContext(ctx, `PRAGMA foreign_keys`).Scan(&enabled); err != nil || enabled != 1 {
			t.Fatalf("foreign_keys = %d, %v", enabled, err)
		}
		if _, err := conn.ExecContext(ctx, `INSERT INTO members(id,project_id,name) VALUES(?,?,?)`, uuid.NewString(), "missing", "Ana"); err == nil {
			t.Fatal("orphan accepted")
		}
	}
}

func TestMembersSQLiteRejectsDuplicatesAndEmptyNamesWithoutOverwriting(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	p := createMemberProject(t, db, "Carva Demo")
	h := memberHandler(db)
	m := decodeMember(t, memberRequest(h, "POST", p, `{"name":"Ana"}`))
	for _, args := range [][]any{{m.ID, p, "Bruno"}, {uuid.NewString(), p, ""}} {
		if _, err := db.Exec(`INSERT INTO members(id,project_id,name) VALUES(?,?,?)`, args...); err == nil {
			t.Fatal("constraint not enforced")
		}
	}
	if got := listMembers(t, h, p); !reflect.DeepEqual(got, []member.Member{m}) {
		t.Fatal(got)
	}
}

func TestMembersSQLiteAbortedWriteLeavesNoRegistration(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	p := createMemberProject(t, db, "Carva Demo")
	if _, err := db.Exec(`CREATE TRIGGER reject_member BEFORE INSERT ON members BEGIN SELECT RAISE(ABORT, 'controlled failure'); END`); err != nil {
		t.Fatal(err)
	}
	h := memberHandler(db)
	w := memberRequest(h, "POST", p, `{"name":"Ana"}`)
	if w.Code != 500 || strings.Contains(w.Body.String(), "controlled failure") {
		t.Fatalf("failure = %d %s", w.Code, w.Body)
	}
	if got := listMembers(t, h, p); len(got) != 0 {
		t.Fatalf("partial registration: %v", got)
	}
}

func TestMembersMigrationUpgradesExistingDatabase(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	// Reconstruct the previous schema using the unchanged applied migrations.
	if _, err := db.Exec(`DROP TABLE members; DELETE FROM carva_schema_migrations WHERE version >= 6`); err != nil {
		t.Fatal(err)
	}
	p := createMemberProject(t, db, "Existing project")
	if err := sqlite.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	if err := sqlite.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	h := memberHandler(db)
	got := decodeMember(t, memberRequest(h, "POST", p, `{"name":"Ana"}`))
	if got.ProjectID != p {
		t.Fatal(got)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM carva_schema_migrations WHERE version = 6`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("migration ledger: %d, %v", count, err)
	}
	project, err := projects.NewService(sqlite.NewProjectRepository(db)).Get(context.Background(), p)
	if err != nil || project.Name != "Existing project" {
		t.Fatalf("existing project lost: %v %v", project, err)
	}
}

func TestMembersNameStorageMigrationPreservesExistingIdentity(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	p := createMemberProject(t, db, "Carva Demo")
	h := memberHandler(db)
	m := decodeMember(t, memberRequest(h, "POST", p, `{"name":"Ana"}`))
	// The pre-007 table has the same columns; replaying 007 verifies its data copy.
	if _, err := db.Exec(`DELETE FROM carva_schema_migrations WHERE version = 7`); err != nil {
		t.Fatal(err)
	}
	if err := sqlite.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	if got := listMembers(t, h, p); !reflect.DeepEqual(got, []member.Member{m}) {
		t.Fatalf("migration changed identity: %v", got)
	}
}
