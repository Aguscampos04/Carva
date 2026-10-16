package project

import "testing"

func TestNewProject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Project
		wantErr bool
	}{
		{name: "trims name and starts planned", input: " Carva ", want: Project{Name: "Carva", Status: Planned}},
		{name: "rejects blank name", input: "  ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && (got.Name != tt.want.Name || got.Status != tt.want.Status) {
				t.Fatalf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestStartFirstSprint(t *testing.T) {
	tests := []struct {
		name string
		in   Status
		want Status
	}{
		{name: "planned project starts", in: Planned, want: InProgress},
		{name: "in progress project remains", in: InProgress, want: InProgress},
		{name: "finalized project remains", in: Finalized, want: Finalized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Project{Status: tt.in}
			p.StartFirstSprint()
			if p.Status != tt.want {
				t.Fatalf("status = %q, want %q", p.Status, tt.want)
			}
		})
	}
}
