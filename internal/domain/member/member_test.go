package member_test

import (
	"errors"
	"testing"

	"github.com/Aguscampos04/Carva/internal/domain/member"
)

// US-03: BR-02, BDD-02/03/15/17.
func TestMemberName(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"  Ana  Pérez  ", "Ana  Pérez"}, {" a ", "a"},
		{"\t\n\u00a0Éva O'Neill\u00a0\n\t", "Éva O'Neill"},
		{"<b>Ana</b>", "<b>Ana</b>"}, {"Ana\t Pérez", "Ana\t Pérez"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			m, err := member.New("member-id", "project-id", tc.input)
			if err != nil || m.Name != tc.want || m.ID != "member-id" || m.ProjectID != "project-id" {
				t.Fatalf("New = %+v, %v; want name %q and unchanged IDs", m, err, tc.want)
			}
		})
	}
	for _, input := range []string{"", "   ", "\t\n\u00a0"} {
		if _, err := member.New("id", "project", input); !errors.Is(err, member.ErrInvalidName) {
			t.Errorf("name %q: %v", input, err)
		}
	}
}

func TestMemberRequiresIdentityAndProject(t *testing.T) {
	for _, ids := range [][2]string{{"", "project"}, {"id", ""}, {" ", "project"}} {
		if _, err := member.New(ids[0], ids[1], "Ana"); err == nil {
			t.Errorf("accepted identity/project %q", ids)
		}
	}
}
