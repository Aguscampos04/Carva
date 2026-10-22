package integration_test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"testing"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	"github.com/Aguscampos04/Carva/internal/application/sprints"
	backlogdomain "github.com/Aguscampos04/Carva/internal/domain/backlog"
	sprintdomain "github.com/Aguscampos04/Carva/internal/domain/sprint"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
)

func TestCloseSprintTransfersOnlyIncompleteStoriesAndKeepsHistory(t *testing.T) {
	fixture := newSprintClosureFixture(t)
	incomplete := fixture.createStory(t, "Incomplete")
	completed := fixture.createStory(t, "Completed")
	fixture.assign(t, fixture.current, incomplete.ID)
	fixture.assign(t, fixture.current, completed.ID)
	fixture.assign(t, fixture.next, incomplete.ID)
	if _, err := fixture.service.CompleteBacklogItem(context.Background(), fixture.projectID, fixture.current.ID, completed.ID, true); err != nil {
		t.Fatal(err)
	}

	before, err := fixture.repo.ListPlannedBacklogItems(context.Background(), fixture.projectID, fixture.current.ID)
	if err != nil {
		t.Fatal(err)
	}
	closed, err := fixture.repo.Close(context.Background(), fixture.projectID, fixture.current.ID)
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !closed.Closed {
		t.Fatalf("closed Sprint = %+v, want closed", closed)
	}
	after, err := fixture.repo.ListPlannedBacklogItems(context.Background(), fixture.projectID, fixture.current.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("closed Sprint history changed: before %+v, after %+v", before, after)
	}
	destination, err := fixture.repo.ListPlannedBacklogItems(context.Background(), fixture.projectID, fixture.next.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(destination) != 1 || destination[0].BacklogItemID != incomplete.ID || destination[0].AcceptanceCriteriaVerified {
		t.Fatalf("destination stories = %+v, want one incomplete story without a duplicate", destination)
	}
}

func TestCloseSprintRequiresDestinationOnlyForIncompleteStories(t *testing.T) {
	t.Run("incomplete story without destination does not partially close", func(t *testing.T) {
		fixture := newSprintClosureFixture(t)
		item := fixture.createStory(t, "Incomplete")
		fixture.assign(t, fixture.current, item.ID)
		if _, err := fixture.repo.Close(context.Background(), fixture.projectID, fixture.current.ID); !errors.Is(err, sprintdomain.ErrNextSprintRequired) {
			t.Fatalf("Close() error = %v, want ErrNextSprintRequired", err)
		}
		got, err := fixture.repo.Get(context.Background(), fixture.projectID, fixture.current.ID)
		if err != nil || got.Closed {
			t.Fatalf("Sprint after rejected close = %+v, %v; want open", got, err)
		}
	})

	t.Run("empty Sprint closes without destination", func(t *testing.T) {
		fixture := newSprintClosureFixture(t)
		fixture.next = sprintdomain.Sprint{}
		if _, err := fixture.db.Exec(`DELETE FROM sprints WHERE id = ?`, fixture.current.ID+1); err != nil {
			t.Fatal(err)
		}
		closed, err := fixture.repo.Close(context.Background(), fixture.projectID, fixture.current.ID)
		if err != nil || !closed.Closed {
			t.Fatalf("Close() = %+v, %v; want closed Sprint", closed, err)
		}
	})
}

func TestCloseSprintRollsBackTransferWhenPersistenceFails(t *testing.T) {
	fixture := newSprintClosureFixture(t)
	item := fixture.createStory(t, "Incomplete")
	fixture.assign(t, fixture.current, item.ID)
	if _, err := fixture.db.Exec(`
		CREATE TRIGGER reject_rollover BEFORE INSERT ON sprint_backlog_items
		WHEN NEW.sprint_id = ` + strconv.FormatInt(fixture.next.ID, 10) + ` BEGIN SELECT RAISE(ABORT, 'controlled rollover failure'); END`); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.repo.Close(context.Background(), fixture.projectID, fixture.current.ID); err == nil {
		t.Fatal("Close() succeeded despite destination insert failure")
	}
	current, err := fixture.repo.Get(context.Background(), fixture.projectID, fixture.current.ID)
	if err != nil || current.Closed {
		t.Fatalf("source Sprint after rollback = %+v, %v; want open", current, err)
	}
	destination, err := fixture.repo.ListPlannedBacklogItems(context.Background(), fixture.projectID, fixture.next.ID)
	if err != nil || len(destination) != 0 {
		t.Fatalf("destination after rollback = %+v, %v; want no partial transfer", destination, err)
	}
}

type sprintClosureFixture struct {
	db        *sql.DB
	repo      *sqlite.SprintRepository
	service   sprints.Service
	projectID string
	current   sprintdomain.Sprint
	next      sprintdomain.Sprint
	backlog   *appbacklog.Service
}

func newSprintClosureFixture(t *testing.T) *sprintClosureFixture {
	t.Helper()
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(context.Background(), "US-10 project")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	sprintRepository := sqlite.NewSprintRepository(db)
	sprintService := sprints.NewService(sprintRepository, backlogRepository)
	current, err := sprintService.Create(context.Background(), projectValue.ID, "Current Sprint")
	if err != nil {
		t.Fatal(err)
	}
	next, err := sprintService.Create(context.Background(), projectValue.ID, "Next Sprint")
	if err != nil {
		t.Fatal(err)
	}
	return &sprintClosureFixture{db: db, repo: sprintRepository, service: sprintService, projectID: projectValue.ID, current: current, next: next,
		backlog: appbacklog.NewService(backlogRepository)}
}

func (fixture *sprintClosureFixture) createStory(t *testing.T, title string) backlogdomain.Item {
	t.Helper()
	item, err := fixture.backlog.Create(context.Background(), fixture.projectID, backlogdomain.Input{
		Title: title, Description: "Description", Priority: backlogdomain.PriorityHigh, AcceptanceCriteria: []string{"Done"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func (fixture *sprintClosureFixture) assign(t *testing.T, target sprintdomain.Sprint, itemID string) {
	t.Helper()
	if _, err := fixture.service.AssignBacklogItem(context.Background(), fixture.projectID, target.ID, itemID); err != nil {
		t.Fatal(err)
	}
}
