// Package planningpoker defines Planning Poker rounds and votes.
package planningpoker

import (
	"errors"
	"strings"
)

var (
	ErrInvalidMemberID      = errors.New("member ID must not be blank")
	ErrInvalidStoryID       = errors.New("story ID must not be blank")
	ErrInvalidRoundID       = errors.New("round ID must not be blank")
	ErrInvalidProjectID     = errors.New("project ID must not be blank")
	ErrInvalidStoryPoints   = errors.New("Story Points must be a non-negative integer")
	ErrNoParticipants       = errors.New("a Planning Poker round requires at least one participant")
	ErrInvalidRoundNumber   = errors.New("round number must be positive")
	ErrDuplicateParticipant = errors.New("Planning Poker round contains a duplicate participant")
	ErrInvalidVoteScope     = errors.New("vote does not belong to this round and story")
	ErrNotParticipant       = errors.New("member is not a participant in this round")
	ErrRoundNotOpen         = errors.New("Planning Poker round is not open")
	ErrRoundNotFound        = errors.New("Planning Poker round not found")
	ErrNotFound             = errors.New("Planning Poker resource not found")
	ErrRoundExists          = errors.New("an initial Planning Poker round already exists for this story")
)

// Vote is one member's Story Points selection for one story and round.
type Vote struct {
	MemberID    string `json:"member_id"`
	StoryID     string `json:"story_id"`
	RoundID     string `json:"round_id"`
	StoryPoints int    `json:"story_points"`
}

// NewVote creates a vote after validating its identity and non-negative value.
func NewVote(memberID, storyID, roundID string, storyPoints int) (Vote, error) {
	if strings.TrimSpace(memberID) == "" {
		return Vote{}, ErrInvalidMemberID
	}
	if strings.TrimSpace(storyID) == "" {
		return Vote{}, ErrInvalidStoryID
	}
	if strings.TrimSpace(roundID) == "" {
		return Vote{}, ErrInvalidRoundID
	}
	if storyPoints < 0 {
		return Vote{}, ErrInvalidStoryPoints
	}
	return Vote{MemberID: memberID, StoryID: storyID, RoundID: roundID, StoryPoints: storyPoints}, nil
}
