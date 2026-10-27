package defects

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Aguscampos04/Carva/internal/domain/defect"
)

type updateRepository struct {
	Repository
	err    error
	called bool
}

func (r *updateRepository) Update(_ context.Context, projectID, storyID, id string, change domain.Patch) (domain.Defect, error) {
	r.called = true
	return domain.Defect{ID: id, ProjectID: projectID, BacklogItemID: storyID, Status: *change.Status}, r.err
}
func TestUpdateValidatesBeforePersistenceAndPropagatesErrors(t *testing.T) {
	status := domain.StatusResolved
	for _, want := range []error{nil, domain.ErrNotFound, errors.New("storage")} {
		r := &updateRepository{err: want}
		got, err := NewService(r).Update(context.Background(), "p", "a", "d", domain.Patch{Status: &status})
		if !errors.Is(err, want) || !r.called {
			t.Fatalf("%+v %v called=%v", got, err, r.called)
		}
		if want == nil && (got.ID != "d" || got.Status != status) {
			t.Fatal(got)
		}
	}
	r := &updateRepository{}
	if _, err := NewService(r).Update(context.Background(), "p", "a", "d", domain.Patch{}); !errors.Is(err, domain.ErrEmptyPatch) || r.called {
		t.Fatalf("%v called=%v", err, r.called)
	}
}
