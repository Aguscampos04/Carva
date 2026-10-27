package integration_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	appdefect "github.com/Aguscampos04/Carva/internal/application/defects"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func defectHandler(db *sql.DB) http.Handler {
	backlogRepository := sqlite.NewBacklogRepository(db)
	base := httpapi.NewHandler(
		projects.NewService(sqlite.NewProjectRepository(db)),
		appbacklog.NewService(backlogRepository),
		sprints.NewService(sqlite.NewSprintRepository(db), backlogRepository),
	)
	return httpapi.WithDefects(base, appdefect.NewService(sqlite.NewDefectRepository(db)))
}

func defectRequest(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestUS18RegisterAndListDefectsWithOptionalSprints(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	seedEffort(t, db)
	result, err := db.Exec(`INSERT INTO sprints(project_id,goal) VALUES('p','Sprint de Carva'),('q','Sprint ajeno')`)
	if err != nil {
		t.Fatal(err)
	}
	firstSprint, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	secondSprint := firstSprint + 1
	handler := defectHandler(db)

	withoutSprint := `{"description":"  No guarda cambios  ","severity":"Alta","status":"Abierto"}`
	response := defectRequest(handler, "POST", "/api/projects/p/backlog/a/defects", withoutSprint)
	expectCode(t, response, http.StatusCreated)
	created := readObject(t, response)
	if created["description"] != "No guarda cambios" || created["backlog_item_id"] != "a" || created["detection_sprint_id"] != nil || created["resolution_sprint_id"] != nil {
		t.Fatalf("unexpected defect without Sprints: %#v", created)
	}

	withDetection := `{"description":"Error al exportar","severity":"Media","status":"En progreso","detection_sprint_id":` + sprintJSON(firstSprint) + `}`
	response = defectRequest(handler, "POST", "/api/projects/p/backlog/a/defects", withDetection)
	expectCode(t, response, http.StatusCreated)
	registered := readObject(t, response)
	if registered["detection_sprint_id"] != float64(firstSprint) || registered["resolution_sprint_id"] != nil {
		t.Fatalf("detection Sprint must be retained without a resolution Sprint: %#v", registered)
	}

	listed := defectRequest(handler, "GET", "/api/projects/p/backlog/a/defects", "")
	expectCode(t, listed, http.StatusOK)
	var defects []map[string]any
	if err := json.Unmarshal(listed.Body.Bytes(), &defects); err != nil || len(defects) != 2 {
		t.Fatalf("defects = %#v, err = %v", defects, err)
	}

	wrongProjectSprint := `{"description":"Foreign Sprint","severity":"Baja","status":"Abierto","detection_sprint_id":` + sprintJSON(secondSprint) + `}`
	expectCode(t, defectRequest(handler, "POST", "/api/projects/p/backlog/a/defects", wrongProjectSprint), http.StatusNotFound)
	for _, body := range []string{
		`{"description":" ","severity":"Alta","status":"Abierto"}`,
		`{"description":"Bug","severity":"Grave","status":"Abierto"}`,
		`{"description":"Bug","severity":"Alta","status":"Pendiente"}`,
	} {
		expectCode(t, defectRequest(handler, "POST", "/api/projects/p/backlog/a/defects", body), http.StatusBadRequest)
	}
}

func sprintJSON(id int64) string { return strconv.FormatInt(id, 10) }
