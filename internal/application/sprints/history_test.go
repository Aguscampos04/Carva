package sprints

import (
	"context"
	"errors"
	"testing"

	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
)

type historyRepository struct {
	Repository
	projectErr, listErr error
	called              bool
}

func (r *historyRepository) GetProject(context.Context, string) (project.Project, error) {
	return project.Project{ID: "project"}, r.projectErr
}
func (r *historyRepository) ListClosed(context.Context, string) ([]sprint.Sprint, error) {
	r.called = true
	return []sprint.Sprint{{ID: 1, ProjectID: "project", Goal: "History", Closed: true}}, r.listErr
}

func TestListClosedValidatesProjectAndPropagatesErrors(t *testing.T) {
	failure := errors.New("storage unavailable")
	for _, tt := range []struct {
		name, id                  string
		projectErr, listErr, want error
		called                    bool
	}{
		{name: "valid", id: "project", called: true},
		{name: "empty", want: project.ErrNotFound},
		{name: "missing", id: "missing", projectErr: project.ErrNotFound, want: project.ErrNotFound},
		{name: "project storage failure", id: "project", projectErr: failure, want: failure},
		{name: "list storage failure", id: "project", listErr: failure, want: failure, called: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := &historyRepository{projectErr: tt.projectErr, listErr: tt.listErr}
			got, err := NewService(repo, nil).ListClosed(context.Background(), tt.id)
			if !errors.Is(err, tt.want) || repo.called != tt.called {
				t.Fatalf("got %+v, %v, called=%v", got, err, repo.called)
			}
			if tt.want == nil && (len(got) != 1 || !got[0].Closed) {
				t.Fatal(got)
			}
		})
	}
}
