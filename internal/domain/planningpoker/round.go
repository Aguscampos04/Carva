package planningpoker

import (
	"slices"
	"strings"
)

// RoundStatus describes whether votes may still be changed.
type RoundStatus string

const (
	RoundOpen      RoundStatus = "open"
	RoundFinalized RoundStatus = "finalized"
)

// Round snapshots the project members participating when voting starts.
type Round struct {
	ID           string      `json:"id"`
	ProjectID    string      `json:"project_id"`
	StoryID      string      `json:"story_id"`
	Number       int         `json:"number"`
	Status       RoundStatus `json:"status"`
	Participants []string    `json:"participants"`
}

// NewRound opens an estimation round for its project members.
func NewRound(id, projectID, storyID string, number int, participants []string) (Round, error) {
	if strings.TrimSpace(id) == "" {
		return Round{}, ErrInvalidRoundID
	}
	if strings.TrimSpace(projectID) == "" {
		return Round{}, ErrInvalidProjectID
	}
	if strings.TrimSpace(storyID) == "" {
		return Round{}, ErrInvalidStoryID
	}
	if number < 1 {
		return Round{}, ErrInvalidRoundNumber
	}
	if len(participants) == 0 {
		return Round{}, ErrNoParticipants
	}
	seen := make(map[string]struct{}, len(participants))
	for _, participant := range participants {
		if strings.TrimSpace(participant) == "" {
			return Round{}, ErrInvalidMemberID
		}
		if _, exists := seen[participant]; exists {
			return Round{}, ErrDuplicateParticipant
		}
		seen[participant] = struct{}{}
	}
	return Round{ID: id, ProjectID: projectID, StoryID: storyID, Number: number, Status: RoundOpen, Participants: slices.Clone(participants)}, nil
}

// ValidateVote ensures the vote belongs to this open round and a participant.
func (round Round) ValidateVote(vote Vote) error {
	if round.Status != RoundOpen {
		return ErrRoundNotOpen
	}
	if vote.RoundID != round.ID || vote.StoryID != round.StoryID {
		return ErrInvalidVoteScope
	}
	if !round.hasParticipant(vote.MemberID) {
		return ErrNotParticipant
	}
	return nil
}

// VisibleVotes returns only the requesting participant's vote while open.
// Finalized rounds reveal the full set together.
func (round Round) VisibleVotes(memberID string, votes []Vote) ([]Vote, error) {
	if !round.hasParticipant(memberID) {
		return nil, ErrNotParticipant
	}
	if round.Status == RoundFinalized {
		return slices.Clone(votes), nil
	}
	visible := make([]Vote, 0, 1)
	for _, vote := range votes {
		if vote.MemberID == memberID && vote.RoundID == round.ID && vote.StoryID == round.StoryID {
			visible = append(visible, vote)
		}
	}
	return visible, nil
}

func (round Round) hasParticipant(memberID string) bool {
	return slices.Contains(round.Participants, memberID)
}
