package planningpoker

import (
	"errors"
	"testing"
)

func TestRoundFinalizeRequiresEveryParticipantToVote(t *testing.T) {
	round := testRound(t)
	result, err := round.Finalize([]Vote{{MemberID: "member-1", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 5}})
	if !errors.Is(err, ErrVotesPending) {
		t.Fatalf("Finalize() error = %v, want %v", err, ErrVotesPending)
	}
	if result.Round.Status != "" || len(result.Votes) != 0 {
		t.Fatalf("Finalize() returned partial result %+v", result)
	}
}

func TestRoundFinalizeRevealsAllVotesAndDetectsDifferences(t *testing.T) {
	round := testRound(t)
	votes := []Vote{
		{MemberID: "member-1", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 3},
		{MemberID: "member-2", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 8},
	}

	result, err := round.Finalize(votes)
	if err != nil {
		t.Fatalf("Finalize() error = %v", err)
	}
	if result.Round.Status != RoundFinalized {
		t.Fatalf("result status = %q, want %q", result.Round.Status, RoundFinalized)
	}
	if len(result.Votes) != 2 {
		t.Fatalf("revealed votes = %d, want 2", len(result.Votes))
	}
	if !result.HasDifferences {
		t.Fatal("HasDifferences = false, want true")
	}
}

func TestRoundFinalizedResultReportsConsensus(t *testing.T) {
	round := testRound(t)
	votes := []Vote{
		{MemberID: "member-1", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 5},
		{MemberID: "member-2", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 5},
	}

	finalized, err := round.Finalize(votes)
	if err != nil {
		t.Fatalf("Finalize() error = %v", err)
	}
	result, err := finalized.Round.Result(votes)
	if err != nil {
		t.Fatalf("Result() error = %v", err)
	}
	if result.HasDifferences {
		t.Fatal("HasDifferences = true, want false")
	}
	if len(result.Votes) != 2 {
		t.Fatalf("result votes = %d, want 2", len(result.Votes))
	}
}

func TestRoundResultAgreedStoryPointsRequiresFinalizedUnanimousVotes(t *testing.T) {
	round := testRound(t)
	finalizedResult := func(t *testing.T, points ...int) RoundResult {
		t.Helper()
		votes := make([]Vote, len(points))
		for i, value := range points {
			votes[i] = Vote{MemberID: round.Participants[i], StoryID: round.StoryID, RoundID: round.ID, StoryPoints: value}
		}
		finalized, err := round.Finalize(votes)
		if err != nil {
			t.Fatalf("Finalize() error = %v", err)
		}
		return finalized
	}
	for _, test := range []struct {
		name      string
		result    RoundResult
		want      int
		wantError error
	}{
		{
			name:   "unanimous zero is a valid agreement",
			result: finalizedResult(t, 0, 0),
		},
		{
			name:   "unanimous nonzero agreement",
			result: finalizedResult(t, 5, 5),
			want:   5,
		},
		{
			name:      "different votes have no agreement",
			result:    finalizedResult(t, 3, 8),
			wantError: ErrNoConsensus,
		},
		{
			name:      "open round cannot be agreed",
			result:    RoundResult{Round: round},
			wantError: ErrRoundNotFinalized,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.result.AgreedStoryPoints()
			if !errors.Is(err, test.wantError) {
				t.Fatalf("AgreedStoryPoints() error = %v, want %v", err, test.wantError)
			}
			if err == nil && got != test.want {
				t.Fatalf("AgreedStoryPoints() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestRoundResultCannotBeReadBeforeFinalization(t *testing.T) {
	round := testRound(t)
	_, err := round.Result(nil)
	if !errors.Is(err, ErrRoundNotFinalized) {
		t.Fatalf("Result() error = %v, want %v", err, ErrRoundNotFinalized)
	}
}

func TestRoundFinalizeRejectsVoteFromNonParticipant(t *testing.T) {
	round := testRound(t)
	votes := []Vote{
		{MemberID: "member-1", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 5},
		{MemberID: "outsider", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 5},
	}
	_, err := round.Finalize(votes)
	if !errors.Is(err, ErrNotParticipant) {
		t.Fatalf("Finalize() error = %v, want %v", err, ErrNotParticipant)
	}
}

func TestRoundFinalizeRejectsDuplicateVotesForOneParticipant(t *testing.T) {
	round := testRound(t)
	votes := []Vote{
		{MemberID: "member-1", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 3},
		{MemberID: "member-1", StoryID: round.StoryID, RoundID: round.ID, StoryPoints: 8},
	}
	_, err := round.Finalize(votes)
	if !errors.Is(err, ErrInvalidVoteSet) {
		t.Fatalf("Finalize() error = %v, want %v", err, ErrInvalidVoteSet)
	}
}

func testRound(t *testing.T) Round {
	t.Helper()
	round, err := NewRound("round-1", "project-1", "story-1", 1, []string{"member-1", "member-2"})
	if err != nil {
		t.Fatalf("NewRound() error = %v", err)
	}
	return round
}
