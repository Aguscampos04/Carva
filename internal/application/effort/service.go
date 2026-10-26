package effort

import (
	"context"
	domain "github.com/Aguscampos04/Carva/internal/domain/effort"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/google/uuid"
	"time"
	_ "time/tzdata"
)

type SprintOption struct {
	ID     int64 `json:"id"`
	Closed bool  `json:"closed"`
}
type FormContext struct {
	Today   project.Date   `json:"today"`
	Sprints []SprintOption `json:"sprints"`
}
type Repository interface {
	Save(context.Context, domain.Entry, bool, func(project.Project, hours.Value) error) (domain.Entry, error)
	Get(context.Context, string, string) (domain.Entry, error)
	List(context.Context, string, string) ([]domain.Entry, error)
	Delete(context.Context, string, string, string) error
	Sprints(context.Context, string, string) ([]SprintOption, error)
}
type Service struct {
	repository Repository
	clock      func() time.Time
}

func NewService(repository Repository, clock func() time.Time) *Service {
	return &Service{repository, clock}
}
func (s *Service) today() project.Date {
	zone, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	return project.Date(s.clock().In(zone).Format(time.DateOnly))
}
func (s *Service) Save(ctx context.Context, pid, id string, in domain.Input) (domain.Entry, error) {
	if err := in.Validate(); err != nil {
		return domain.Entry{}, err
	}
	create := id == ""
	if create {
		id = uuid.NewString()
	}
	return s.repository.Save(ctx, domain.Entry{ID: id, ProjectID: pid, Input: in}, create, func(p project.Project, other hours.Value) error {
		if err := domain.ValidateDate(p, s.today(), in.WorkDate); err != nil {
			return err
		}
		return domain.ValidateDaily(other, in.Hours)
	})
}
func (s *Service) Get(ctx context.Context, pid, id string) (domain.Entry, error) {
	return s.repository.Get(ctx, pid, id)
}
func (s *Service) List(ctx context.Context, pid, story string) ([]domain.Entry, error) {
	return s.repository.List(ctx, pid, story)
}
func (s *Service) Delete(ctx context.Context, pid, id string) error {
	return s.repository.Delete(ctx, pid, id, s.clock().UTC().Format(time.RFC3339Nano))
}
func (s *Service) Context(ctx context.Context, pid, story string) (FormContext, error) {
	options, err := s.repository.Sprints(ctx, pid, story)
	return FormContext{Today: s.today(), Sprints: options}, err
}
