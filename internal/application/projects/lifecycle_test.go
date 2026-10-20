package projects

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Aguscampos04/Carva/internal/domain/project"
)

// BDD-13/15: dates are computed on the server in Buenos Aires, once per
// effective transition, and not recomputed by an already-finalized retry.
func TestLifecycleClockBoundaryAndRepeatedFinalization(t *testing.T) {
	for _, test := range []struct{ instant, want string }{
		{"2026-10-21T02:30:00Z", "2026-10-20"},
		{"2026-10-21T03:00:00Z", "2026-10-21"},
	} {
		t.Run(test.instant, func(t *testing.T) {
			repo := newMemoryRepository()
			calls := 0
			service := NewServiceWithClock(repo, func() time.Time { calls++; now, _ := time.Parse(time.RFC3339, test.instant); return now })
			p, err := service.Create(context.Background(), "Demo")
			if err != nil {
				t.Fatal(err)
			}
			final, err := service.Finalize(context.Background(), p.ID)
			if err != nil || string(final.EndDate) != test.want {
				t.Fatalf("%+v %v", final, err)
			}
			again, err := service.Finalize(context.Background(), p.ID)
			if err != nil || again != final || calls != 1 {
				t.Fatalf("repeated: %+v %v calls=%d", again, err, calls)
			}
		})
	}
}

func TestLifecycleApplicationDoesNotPersistInvalidChanges(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo)
	p, err := service.Create(context.Background(), "Demo")
	if err != nil {
		t.Fatal(err)
	}
	blank := " "
	if _, err = service.Update(context.Background(), p.ID, project.Patch{Name: &blank}); !errors.Is(err, project.ErrInvalidName) {
		t.Fatal(err)
	}
	if got := repo.projects[p.ID]; got != p {
		t.Fatal(got)
	}
	name := "Changed"
	failure := errors.New("write failed")
	repo.saveErr = failure
	if _, err = service.Update(context.Background(), p.ID, project.Patch{Name: &name}); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	if got := repo.projects[p.ID]; got != p {
		t.Fatal("partial mutation")
	}
	if _, err = service.Finalize(context.Background(), "missing"); !errors.Is(err, project.ErrNotFound) {
		t.Fatal(err)
	}
}

// BDD-11: recording/passing a calendar day and consulting never changes state.
func TestProjectDateDoesNotDriveState(t *testing.T) {
	repo := newMemoryRepository()
	now := time.Date(2026, 10, 20, 15, 0, 0, 0, time.UTC)
	clockCalls := 0
	service := NewServiceWithClock(repo, func() time.Time { clockCalls++; return now })
	p, err := service.Create(context.Background(), "Demo")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.Update(context.Background(), p.ID, project.Patch{HasStartDate: true, StartDate: project.Date("2026-10-20")}); err != nil {
		t.Fatal(err)
	}
	now = now.AddDate(0, 0, 1)
	got, err := service.Get(context.Background(), p.ID)
	if err != nil || got.Status != project.Planned || got.EndDate != "" || clockCalls != 0 {
		t.Fatalf("%+v %v clock calls %d", got, err, clockCalls)
	}
}
