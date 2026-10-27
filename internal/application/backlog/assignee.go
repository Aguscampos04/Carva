package backlog

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidAssignee = errors.New("responsible member must belong to the same project")

type Assignee struct {
	MemberID *string `json:"member_id"`
}
type AssigneeRepository interface {
	Get(context.Context, string, string) (Assignee, error)
	Set(context.Context, string, string, *string) error
}
type AssigneeService struct{ repository AssigneeRepository }

func NewAssigneeService(r AssigneeRepository) *AssigneeService {
	return &AssigneeService{repository: r}
}
func (s *AssigneeService) Get(ctx context.Context, pid, id string) (Assignee, error) {
	return s.repository.Get(ctx, pid, id)
}
func (s *AssigneeService) Set(ctx context.Context, pid, id string, memberID *string) (Assignee, error) {
	if memberID != nil {
		if _, err := uuid.Parse(*memberID); err != nil {
			return Assignee{}, ErrInvalidAssignee
		}
	}
	if err := s.repository.Set(ctx, pid, id, memberID); err != nil {
		return Assignee{}, err
	}
	return Assignee{MemberID: memberID}, nil
}
