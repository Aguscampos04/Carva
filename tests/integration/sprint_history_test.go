package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	appSprints "github.com/Aguscampos04/Carva/internal/application/sprints"
	backlog "github.com/Aguscampos04/Carva/internal/domain/backlog"
	sprint "github.com/Aguscampos04/Carva/internal/domain/sprint"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

// BDD-01/04/05: list only this project's closed Sprints, in creation order.
func TestUS11ClosedSprintListing(t *testing.T) {
	f := newSprintClosureFixture(t)
	path := "/api/projects/" + f.projectID + "/sprints/closed"
	empty := performJSONRequest(f.handler(), http.MethodGet, path, nil)
	if empty.Code != 200 || empty.Body.String() != "[]\n" {
		t.Fatalf("empty = %d %s", empty.Code, empty.Body)
	}
	for _, id := range []int64{f.next.ID, f.current.ID} {
		if _, err := f.service.Close(context.Background(), f.projectID, id); err != nil {
			t.Fatal(err)
		}
	}
	foreign := newSprintClosureFixture(t) // independent database cannot contaminate this result
	_ = foreign
	otherProject := createProject(t, f.handler(), "Other project")
	if _, err := f.service.Create(context.Background(), otherProject, "Foreign goal"); err != nil {
		t.Fatal(err)
	}
	listed := performJSONRequest(f.handler(), http.MethodGet, path, nil)
	var got []sprint.Sprint
	if listed.Code != 200 || json.Unmarshal(listed.Body.Bytes(), &got) != nil || len(got) != 2 {
		t.Fatalf("list = %d %s", listed.Code, listed.Body)
	}
	if got[0].ID != f.current.ID || got[1].ID != f.next.ID || !got[0].Closed || got[0].Goal != f.current.Goal {
		t.Fatal(got)
	}
	for _, url := range []string{"/api/projects/missing/sprints/closed", sprintPath(otherProject, f.current.ID), sprintPath(f.projectID, 999999)} {
		response := performJSONRequest(f.handler(), http.MethodGet, url, nil)
		if response.Code != 404 {
			t.Fatalf("%s = %d %s", url, response.Code, response.Body)
		}
	}
	if err := f.db.Close(); err != nil {
		t.Fatal(err)
	}
	response := performJSONRequest(f.handler(), http.MethodGet, path, nil)
	if response.Code != 500 {
		t.Fatalf("closed DB = %d %s", response.Code, response.Body)
	}
}

// AC-01/02, BDD-01/02/03: read persisted snapshots after rollover, later edits,
// later completion and reopening the database; reads themselves perform no writes.
func TestUS11HistorySurvivesTransferAndLaterChanges(t *testing.T) {
	ctx := context.Background()
	f := newSprintClosureFixture(t)
	incomplete := f.createStory(t, "Original incomplete")
	completed := f.createStory(t, "Original completed")
	f.assign(t, f.current, incomplete.ID)
	f.assign(t, f.current, completed.ID)
	if _, err := f.service.CompleteBacklogItem(ctx, f.projectID, f.current.ID, completed.ID, true); err != nil {
		t.Fatal(err)
	}
	if _, err := f.service.Close(ctx, f.projectID, f.current.ID); err != nil {
		t.Fatal(err)
	}
	before, err := f.service.ListPlannedBacklogItems(ctx, f.projectID, f.current.ID)
	if err != nil {
		t.Fatal(err)
	}
	points := 8
	if _, err := f.backlog.Update(ctx, f.projectID, incomplete.ID, backlog.Input{Title: "Changed title", Description: "Changed description", Priority: backlog.PriorityLow, ProvisionalStoryPoints: &points, AcceptanceCriteria: []string{"Changed criterion"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := f.service.CompleteBacklogItem(ctx, f.projectID, f.next.ID, incomplete.ID, true); err != nil {
		t.Fatal(err)
	}
	var dbPath string
	var sequence int
	var dbName string
	if err := f.db.QueryRow("PRAGMA database_list").Scan(&sequence, &dbName, &dbPath); err != nil {
		t.Fatal(err)
	}
	if err := f.db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := sqlite.OpenPath(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	f.db = reopened
	f.repo = sqlite.NewSprintRepository(reopened)
	f.service = appSprints.NewService(f.repo, sqlite.NewBacklogRepository(reopened))
	handler := f.handler()
	for i := 0; i < 2; i++ {
		response := performJSONRequest(handler, http.MethodGet, sprintPath(f.projectID, f.current.ID)+"/backlog-items", nil)
		var after []sprint.PlannedBacklogItem
		if response.Code != 200 || json.Unmarshal(response.Body.Bytes(), &after) != nil || !reflect.DeepEqual(after, before) {
			t.Fatalf("history changed: %d %s; want %+v", response.Code, response.Body, before)
		}
	}
	if len(before) != 2 || before[0].AcceptanceCriteriaVerified || !before[1].AcceptanceCriteriaVerified {
		t.Fatal(before)
	}
	destination, err := f.service.ListPlannedBacklogItems(ctx, f.projectID, f.next.ID)
	if err != nil || len(destination) != 1 || !destination[0].AcceptanceCriteriaVerified {
		t.Fatalf("destination=%+v %v", destination, err)
	}
}
