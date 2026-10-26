package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	appbacklog "github.com/Aguscampos04/Carva/internal/application/backlog"
	"github.com/Aguscampos04/Carva/internal/application/members"
	appplanningpoker "github.com/Aguscampos04/Carva/internal/application/planningpoker"
	"github.com/Aguscampos04/Carva/internal/application/projects"
	domainbacklog "github.com/Aguscampos04/Carva/internal/domain/backlog"
	domainpoker "github.com/Aguscampos04/Carva/internal/domain/planningpoker"
	"github.com/Aguscampos04/Carva/internal/persistence/sqlite"
	"github.com/Aguscampos04/Carva/internal/presentation/httpapi"
)

func TestPlanningPokerVotesAreAttributedAndHiddenUntilFinalization(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	ctx := context.Background()
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(ctx, "Planning Poker")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := appbacklog.NewService(backlogRepository)
	item, err := backlogService.Create(ctx, projectValue.ID, domainbacklog.Input{
		Title: "Login", Description: "A member can sign in", Priority: domainbacklog.PriorityHigh,
		AcceptanceCriteria: []string{"A valid login succeeds"},
	})
	if err != nil {
		t.Fatal(err)
	}
	memberRepository := sqlite.NewMemberRepository(db)
	memberService := members.NewService(memberRepository)
	first, err := memberService.Create(ctx, projectValue.ID, "Ana")
	if err != nil {
		t.Fatal(err)
	}
	second, err := memberService.Create(ctx, projectValue.ID, "Bea")
	if err != nil {
		t.Fatal(err)
	}
	service := appplanningpoker.NewService(sqlite.NewPlanningPokerRepository(db), backlogRepository, memberService)
	handler := httpapi.NewHandlerWithPlanningPoker(projectService, backlogService, nil, service, memberService)
	base := "/api/projects/" + projectValue.ID + "/backlog/" + item.ID + "/estimation/rounds"

	started := performJSONRequest(handler, http.MethodPost, base, nil)
	if started.Code != http.StatusCreated {
		t.Fatalf("start round = %d %s", started.Code, started.Body)
	}
	var round domainpoker.Round
	if err := json.Unmarshal(started.Body.Bytes(), &round); err != nil {
		t.Fatal(err)
	}
	if round.Number != 1 || round.Status != domainpoker.RoundOpen || len(round.Participants) != 2 {
		t.Fatalf("initial round = %+v, want round 1, open, with both registered members", round)
	}
	if duplicate := performJSONRequest(handler, http.MethodPost, base, nil); duplicate.Code != http.StatusConflict {
		t.Fatalf("second initial round = %d %s, want conflict", duplicate.Code, duplicate.Body)
	}

	cast := func(memberID string, points int) *httptest.ResponseRecorder {
		t.Helper()
		path := base + "/" + round.ID + "/votes/" + memberID
		return performJSONRequest(handler, http.MethodPut, path, map[string]any{"story_points": points})
	}
	readVotes := func(memberID string) []domainpoker.Vote {
		t.Helper()
		path := base + "/" + round.ID + "/votes?member_id=" + url.QueryEscape(memberID)
		response := performJSONRequest(handler, http.MethodGet, path, nil)
		if response.Code != http.StatusOK {
			t.Fatalf("read votes = %d %s", response.Code, response.Body)
		}
		var votes []domainpoker.Vote
		if err := json.Unmarshal(response.Body.Bytes(), &votes); err != nil {
			t.Fatal(err)
		}
		return votes
	}
	if response := cast(first.ID, 3); response.Code != http.StatusOK {
		t.Fatalf("first vote = %d %s", response.Code, response.Body)
	}
	if visible := readVotes(first.ID); len(visible) != 1 || visible[0].MemberID != first.ID {
		t.Fatalf("first participant saw %#v before everyone voted, want only their vote", visible)
	}
	if visible := readVotes(second.ID); len(visible) != 0 {
		t.Fatalf("second participant saw %#v before voting, want no votes", visible)
	}
	finalizePath := base + "/" + round.ID + "/finalize?member_id=" + url.QueryEscape(first.ID)
	if response := performJSONRequest(handler, http.MethodPost, finalizePath, nil); response.Code != http.StatusConflict {
		t.Fatalf("finalize with pending vote = %d %s, want conflict", response.Code, response.Body)
	}
	resultPath := base + "/" + round.ID + "/result?member_id=" + url.QueryEscape(first.ID)
	if response := performJSONRequest(handler, http.MethodGet, resultPath, nil); response.Code != http.StatusConflict {
		t.Fatalf("read result before finalization = %d %s, want conflict", response.Code, response.Body)
	}
	if response := cast(second.ID, 8); response.Code != http.StatusOK {
		t.Fatalf("second vote = %d %s", response.Code, response.Body)
	}
	if visible := readVotes(first.ID); len(visible) != 1 || visible[0].MemberID != first.ID || visible[0].StoryPoints != 3 {
		t.Fatalf("first participant saw %#v, want only their own vote", visible)
	}
	if visible := readVotes(second.ID); len(visible) != 1 || visible[0].MemberID != second.ID || visible[0].StoryPoints != 8 {
		t.Fatalf("second participant saw %#v, want only their own vote", visible)
	}
	if response := cast(first.ID, 5); response.Code != http.StatusOK {
		t.Fatalf("replace vote = %d %s", response.Code, response.Body)
	}
	if visible := readVotes(first.ID); len(visible) != 1 || visible[0].StoryPoints != 5 {
		t.Fatalf("replaced vote = %#v, want 5", visible)
	}
	if response := cast(first.ID, -1); response.Code != http.StatusBadRequest {
		t.Fatalf("negative vote = %d %s, want bad request", response.Code, response.Body)
	}

	otherProjectID := createMemberProject(t, db, "Other project")
	outsider, err := memberService.Create(ctx, otherProjectID, "Outsider")
	if err != nil {
		t.Fatal(err)
	}
	if response := cast(outsider.ID, 2); response.Code != http.StatusBadRequest {
		t.Fatalf("outsider vote = %d %s, want bad request", response.Code, response.Body)
	}

	stored, err := backlogService.List(ctx, projectValue.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(stored) != 1 || stored[0].StoryPoints != nil {
		t.Fatalf("votes changed definitive Story Points: %+v", stored)
	}

	finalized := performJSONRequest(handler, http.MethodPost, finalizePath, nil)
	if finalized.Code != http.StatusOK {
		t.Fatalf("finalize complete round = %d %s", finalized.Code, finalized.Body)
	}
	var result domainpoker.RoundResult
	if err := json.Unmarshal(finalized.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Round.Status != domainpoker.RoundFinalized || len(result.Votes) != 2 || !result.HasDifferences {
		t.Fatalf("finalized result = %+v, want both votes and differences", result)
	}
	secondResultPath := base + "/" + round.ID + "/result?member_id=" + url.QueryEscape(second.ID)
	secondResult := performJSONRequest(handler, http.MethodGet, secondResultPath, nil)
	if secondResult.Code != http.StatusOK {
		t.Fatalf("second participant read revealed result = %d %s", secondResult.Code, secondResult.Body)
	}
	if response := cast(first.ID, 13); response.Code != http.StatusConflict {
		t.Fatalf("vote after finalization = %d %s, want conflict", response.Code, response.Body)
	}
	if response := performJSONRequest(handler, http.MethodPost, finalizePath, nil); response.Code != http.StatusConflict {
		t.Fatalf("finalize already-finalized round = %d %s, want conflict", response.Code, response.Body)
	}

	startedNext := performJSONRequest(handler, http.MethodPost, base, nil)
	if startedNext.Code != http.StatusCreated {
		t.Fatalf("start next round after finalization = %d %s, want created", startedNext.Code, startedNext.Body)
	}
	var nextRound domainpoker.Round
	if err := json.Unmarshal(startedNext.Body.Bytes(), &nextRound); err != nil {
		t.Fatal(err)
	}
	if nextRound.Number != 2 || nextRound.Status != domainpoker.RoundOpen || len(nextRound.Participants) != 2 {
		t.Fatalf("next round = %+v, want round 2 open with both participants", nextRound)
	}
	listResponse := performJSONRequest(handler, http.MethodGet, base, nil)
	var rounds []domainpoker.Round
	if err := json.Unmarshal(listResponse.Body.Bytes(), &rounds); err != nil {
		t.Fatal(err)
	}
	if listResponse.Code != http.StatusOK || len(rounds) != 2 || rounds[0].Number != 1 || rounds[1].Number != 2 {
		t.Fatalf("round history = %d %+v, want rounds 1 and 2", listResponse.Code, rounds)
	}
	if previous := performJSONRequest(handler, http.MethodGet, resultPath, nil); previous.Code != http.StatusOK {
		t.Fatalf("previous round result after starting another = %d %s, want OK", previous.Code, previous.Body)
	}
	secondRoundVotesPath := base + "/" + nextRound.ID + "/votes?member_id=" + url.QueryEscape(first.ID)
	if votesResponse := performJSONRequest(handler, http.MethodGet, secondRoundVotesPath, nil); votesResponse.Code != http.StatusOK || votesResponse.Body.String() != "[]\n" {
		t.Fatalf("new round's votes = %d %s, want empty list", votesResponse.Code, votesResponse.Body)
	}
	secondRoundVotePath := base + "/" + nextRound.ID + "/votes/" + first.ID
	if voteResponse := performJSONRequest(handler, http.MethodPut, secondRoundVotePath, map[string]any{"story_points": 5}); voteResponse.Code != http.StatusOK {
		t.Fatalf("first vote in new round = %d %s, want OK", voteResponse.Code, voteResponse.Body)
	}
	visibleNextVotes := performJSONRequest(handler, http.MethodGet, secondRoundVotesPath, nil)
	var visibleNext []domainpoker.Vote
	if err := json.Unmarshal(visibleNextVotes.Body.Bytes(), &visibleNext); err != nil {
		t.Fatal(err)
	}
	if visibleNextVotes.Code != http.StatusOK || len(visibleNext) != 1 || visibleNext[0].MemberID != first.ID {
		t.Fatalf("new round visible votes = %d %+v, want only the first participant's vote", visibleNextVotes.Code, visibleNext)
	}
	secondMemberVotesPath := base + "/" + nextRound.ID + "/votes?member_id=" + url.QueryEscape(second.ID)
	secondMemberVisible := performJSONRequest(handler, http.MethodGet, secondMemberVotesPath, nil)
	if secondMemberVisible.Code != http.StatusOK || secondMemberVisible.Body.String() != "[]\n" {
		t.Fatalf("second participant saw another member's new-round vote: %d %s, want empty list", secondMemberVisible.Code, secondMemberVisible.Body)
	}
	if duplicate := performJSONRequest(handler, http.MethodPost, base, nil); duplicate.Code != http.StatusConflict {
		t.Fatalf("start another round while round 2 is open = %d %s, want conflict", duplicate.Code, duplicate.Body)
	}
}

func TestPlanningPokerCannotStartWithoutProjectMembers(t *testing.T) {
	_, db := openMigratedProjectDB(t)
	projectService := projects.NewService(sqlite.NewProjectRepository(db))
	projectValue, err := projectService.Create(context.Background(), "Empty team")
	if err != nil {
		t.Fatal(err)
	}
	backlogRepository := sqlite.NewBacklogRepository(db)
	backlogService := appbacklog.NewService(backlogRepository)
	item, err := backlogService.Create(context.Background(), projectValue.ID, domainbacklog.Input{
		Title: "Story", Description: "Description", Priority: domainbacklog.PriorityMedium,
		AcceptanceCriteria: []string{"Criterion"},
	})
	if err != nil {
		t.Fatal(err)
	}
	memberService := members.NewService(sqlite.NewMemberRepository(db))
	service := appplanningpoker.NewService(sqlite.NewPlanningPokerRepository(db), backlogRepository, memberService)
	handler := httpapi.NewHandlerWithPlanningPoker(projectService, backlogService, nil, service, memberService)
	path := "/api/projects/" + projectValue.ID + "/backlog/" + item.ID + "/estimation/rounds"
	response := performJSONRequest(handler, http.MethodPost, path, nil)
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("start without members = %d %s, want unprocessable entity", response.Code, response.Body)
	}
}
