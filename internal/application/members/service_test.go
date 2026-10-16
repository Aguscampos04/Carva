package members_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Aguscampos04/Carva/internal/application/members"
	"github.com/Aguscampos04/Carva/internal/domain/member"
	"github.com/google/uuid"
)

type repository struct {
	exists                     bool
	checkErr, saveErr, listErr error
	items                      []member.Member
	checks                     int
}

func (r *repository) ProjectExists(context.Context, string) (bool, error) {
	r.checks++
	return r.exists, r.checkErr
}
func (r *repository) Create(_ context.Context, m member.Member) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.items = append(r.items, m)
	return nil
}
func (r *repository) ListByProject(_ context.Context, p string) ([]member.Member, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	var result []member.Member
	for _, m := range r.items {
		if m.ProjectID == p {
			result = append(result, m)
		}
	}
	return result, nil
}

// US-03: BDD-01/04/05/06, identities are independent of names.
func TestMembersRegisterAndList(t *testing.T) {
	r := &repository{exists: true}
	s := members.NewService(r)
	ctx := context.Background()
	empty, err := s.List(ctx, "p")
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("empty = %v, %v", empty, err)
	}
	seen := map[string]bool{}
	for _, p := range []string{"p", "p", "other"} {
		m, err := s.Create(ctx, p, "  Ana  Pérez  ")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := uuid.Parse(m.ID); err != nil || seen[m.ID] {
			t.Fatalf("UUID: %+v", m)
		}
		seen[m.ID] = true
		if m.ProjectID != p || m.Name != "Ana  Pérez" {
			t.Fatalf("member = %+v", m)
		}
	}
	list, err := s.List(ctx, "p")
	if err != nil || len(list) != 2 {
		t.Fatalf("list = %v, %v", list, err)
	}
}

func TestMembersErrorsAndValidationPrecedence(t *testing.T) {
	ctx := context.Background()
	failure := errors.New("storage unavailable")
	r := &repository{}
	s := members.NewService(r)
	if _, err := s.Create(ctx, "missing", " "); !errors.Is(err, member.ErrInvalidName) || r.checks != 0 {
		t.Fatalf("validation precedence: %v checks=%d", err, r.checks)
	}
	if _, err := s.Create(ctx, "missing", "Ana"); !errors.Is(err, member.ErrProjectNotFound) {
		t.Fatal(err)
	}
	if _, err := s.List(ctx, "missing"); !errors.Is(err, member.ErrProjectNotFound) {
		t.Fatal(err)
	}
	r.checkErr = failure
	if _, err := s.Create(ctx, "p", "Ana"); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	if _, err := s.List(ctx, "p"); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	r.checkErr = nil
	r.exists = true
	r.saveErr = failure
	if _, err := s.Create(ctx, "p", "Ana"); !errors.Is(err, failure) || len(r.items) != 0 {
		t.Fatalf("save: %v", err)
	}
	r.listErr = failure
	if list, err := s.List(ctx, "p"); !errors.Is(err, failure) || list != nil {
		t.Fatalf("list failure: %v, %v", list, err)
	}
}
