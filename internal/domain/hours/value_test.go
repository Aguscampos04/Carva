package hours

import (
	"encoding/json"
	"testing"
)

// US-16 BDD 01/04/08/19/20: exact centesimal transport, no rounding.
func TestUS16ExactHours(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Value
	}{{"0.01", 1}, {"1", 100}, {"1.50", 150}, {"24", 2400}, {"92233720368547758.07", 9223372036854775807}} {
		var got Value
		if err := json.Unmarshal([]byte(tc.raw), &got); err != nil || got != tc.want {
			t.Fatalf("%s: %v %v", tc.raw, got, err)
		}
		encoded, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		var round Value
		if json.Unmarshal(encoded, &round) != nil || round != got {
			t.Fatal("inexact roundtrip")
		}
	}
	for _, raw := range []string{"-1", "1.001", "1.500", "1e1", "\"1.25\"", "null", "92233720368547758.08", "true", "0.0001"} {
		var got Value
		if json.Unmarshal([]byte(raw), &got) == nil {
			t.Errorf("accepted %s", raw)
		}
	}
}
