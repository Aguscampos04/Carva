package effort

import (
	"context"
	"errors"
	domain "github.com/Aguscampos04/Carva/internal/domain/effort"
	"github.com/Aguscampos04/Carva/internal/domain/hours"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"testing"
	"time"
)

type repositoryDouble struct {
	project project.Project
	other   hours.Value
	last    domain.Entry
	deleted string
}

func (r *repositoryDouble) Save(_ context.Context, e domain.Entry, _ bool, check func(project.Project, hours.Value) error) (domain.Entry, error) {
	r.last = e
	return e, check(r.project, r.other)
}
func (r *repositoryDouble) Get(context.Context, string, string) (domain.Entry, error) {
	return r.last, nil
}
func (r *repositoryDouble) List(context.Context, string, string) ([]domain.Entry, error) {
	return []domain.Entry{r.last}, nil
}
func (r *repositoryDouble) Delete(_ context.Context, _, _, stamp string) error {
	r.deleted = stamp
	return nil
}
func (r *repositoryDouble) Sprints(context.Context, string, string) ([]SprintOption, error) {
	return []SprintOption{}, nil
}

// BDD 13: deterministic midnight boundary in the approved zone, independently of host/browser.
func TestUS16ServerDayAndApplicationRules(t *testing.T) {
	repo := &repositoryDouble{}
	for _, tc := range []struct{ instant, day string }{{"2026-10-26T02:30:00Z", "2026-10-25"}, {"2026-10-26T03:00:00Z", "2026-10-26"}} {
		instant, _ := time.Parse(time.RFC3339, tc.instant)
		s := NewService(repo, func() time.Time { return instant })
		v, err := s.Context(context.Background(), "p", "a")
		if err != nil || string(v.Today) != tc.day || v.Sprints == nil {
			t.Fatal(v, err)
		}
	}
	clock := func() time.Time { return time.Date(2026, 10, 26, 15, 0, 0, 0, time.UTC) }
	s := NewService(repo, clock)
	in := domain.Input{MemberID: "m", StoryID: "a", WorkDate: "2026-10-25", Activity: "Review", Hours: 1}
	e, err := s.Save(context.Background(), "p", "", in)
	if err != nil || e.ID == "" || e.ProjectID != "p" {
		t.Fatal(e, err)
	}
	same, err := s.Save(context.Background(), "p", e.ID, in)
	if err != nil || same.ID != e.ID {
		t.Fatal(same, err)
	}
	repo.other = 2400
	if _, err := s.Save(context.Background(), "p", e.ID, in); !errors.Is(err, domain.ErrDailyLimit) {
		t.Fatal(err)
	}
	repo.project = project.Project{Status: project.Finalized}
	if _, err := s.Save(context.Background(), "p", e.ID, in); !errors.Is(err, domain.ErrUnknownEnd) {
		t.Fatal(err)
	}
	if err := s.Delete(context.Background(), "p", e.ID); err != nil || repo.deleted != "2026-10-26T15:00:00Z" {
		t.Fatal(repo.deleted, err)
	}
	if _, err := s.Save(context.Background(), "p", "", domain.Input{}); !errors.Is(err, domain.ErrInvalid) {
		t.Fatal(err)
	}
	if _, err := s.Get(context.Background(), "p", e.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.List(context.Background(), "p", "a"); err != nil {
		t.Fatal(err)
	}
}
