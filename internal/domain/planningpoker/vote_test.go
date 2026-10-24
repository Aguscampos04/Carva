package planningpoker

import (
	"errors"
	"testing"
)

func TestNewVote(t *testing.T) {
	tests := []struct {
		name      string
		memberID  string
		storyID   string
		roundID   string
		points    int
		wantError error
	}{
		{name: "valid vote", memberID: "member-1", storyID: "story-1", roundID: "round-1", points: 5},
		{name: "zero points", memberID: "member-1", storyID: "story-1", roundID: "round-1", points: 0},
		{name: "missing member", storyID: "story-1", roundID: "round-1", wantError: ErrInvalidMemberID},
		{name: "missing story", memberID: "member-1", roundID: "round-1", wantError: ErrInvalidStoryID},
		{name: "missing round", memberID: "member-1", storyID: "story-1", wantError: ErrInvalidRoundID},
		{name: "negative points", memberID: "member-1", storyID: "story-1", roundID: "round-1", points: -1, wantError: ErrInvalidStoryPoints},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vote, err := NewVote(tt.memberID, tt.storyID, tt.roundID, tt.points)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("NewVote() error = %v, want %v", err, tt.wantError)
			}
			if err != nil {
				return
			}
			if vote.MemberID != tt.memberID || vote.StoryID != tt.storyID || vote.RoundID != tt.roundID || vote.StoryPoints != tt.points {
				t.Fatalf("NewVote() = %#v, does not preserve input", vote)
			}
		})
	}
}

func TestRoundVisibleVotes(t *testing.T) {
	round, err := NewRound("round-1", "project-1", "story-1", 1, []string{"member-1", "member-2"})
	if err != nil {
		t.Fatalf("NewRound() error = %v", err)
	}
	votes := []Vote{
		{MemberID: "member-1", StoryID: "story-1", RoundID: "round-1", StoryPoints: 3},
		{MemberID: "member-2", StoryID: "story-1", RoundID: "round-1", StoryPoints: 8},
	}

	visible, err := round.VisibleVotes("member-1", votes)
	if err != nil {
		t.Fatalf("VisibleVotes() error = %v", err)
	}
	if len(visible) != 1 || visible[0].MemberID != "member-1" {
		t.Fatalf("open round exposed votes %#v, want only member-1's vote", visible)
	}

	round.Status = RoundFinalized
	visible, err = round.VisibleVotes("member-1", votes)
	if err != nil {
		t.Fatalf("VisibleVotes() finalized error = %v", err)
	}
	if len(visible) != 2 {
		t.Fatalf("finalized round exposed %d votes, want all 2", len(visible))
	}
}

func TestNewRoundRequiresDistinctParticipants(t *testing.T) {
	_, err := NewRound("round-1", "project-1", "story-1", 1, []string{"member-1", "member-1"})
	if !errors.Is(err, ErrDuplicateParticipant) {
		t.Fatalf("NewRound() error = %v, want %v", err, ErrDuplicateParticipant)
	}
}

func TestRoundRejectsVoteFromNonParticipant(t *testing.T) {
	round, err := NewRound("round-1", "project-1", "story-1", 1, []string{"member-1"})
	if err != nil {
		t.Fatalf("NewRound() error = %v", err)
	}
	vote, err := NewVote("member-2", "story-1", "round-1", 5)
	if err != nil {
		t.Fatalf("NewVote() error = %v", err)
	}
	if err := round.ValidateVote(vote); !errors.Is(err, ErrNotParticipant) {
		t.Fatalf("ValidateVote() error = %v, want %v", err, ErrNotParticipant)
	}
	round.Status = RoundFinalized
	vote.MemberID = "member-1"
	if err := round.ValidateVote(vote); !errors.Is(err, ErrRoundNotOpen) {
		t.Fatalf("ValidateVote() finalized error = %v, want %v", err, ErrRoundNotOpen)
	}
}
