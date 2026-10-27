package metrics

import (
	"encoding/json"
	"testing"

	"github.com/Aguscampos04/Carva/internal/domain/hours"
)

func TestCalculateHours(t *testing.T) {
	tests := []struct {
		name      string
		estimates []*hours.Value
		actual    hours.Value
		want      HoursMetrics
		wantError error
	}{
		{
			name:      "exact positive deviation",
			estimates: []*hours.Value{hours.Ptr(250), hours.Ptr(325)},
			actual:    625,
			want:      HoursMetrics{EstimatedHours: decimalHours(575), ActualHours: decimalHours(625), Deviation: decimalHours(50)},
		},
		{
			name:      "negative deviation",
			estimates: []*hours.Value{hours.Ptr(800)},
			actual:    550,
			want:      HoursMetrics{EstimatedHours: decimalHours(800), ActualHours: decimalHours(550), Deviation: decimalHours(-250)},
		},
		{
			name:      "missing estimate keeps actual available",
			estimates: []*hours.Value{hours.Ptr(200), nil},
			actual:    200,
			want:      HoursMetrics{EstimatedHours: nil, ActualHours: decimalHours(200), Deviation: nil},
		},
		{
			name:   "empty totals are zero",
			actual: 0,
			want:   HoursMetrics{EstimatedHours: decimalHours(0), ActualHours: decimalHours(0), Deviation: decimalHours(0)},
		},
		{
			name:      "estimate sum overflow is reported",
			estimates: []*hours.Value{hours.Ptr(hours.Value(1<<63 - 1)), hours.Ptr(1)},
			actual:    0,
			wantError: ErrHoursOverflow,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CalculateHours(test.estimates, test.actual)
			if err != test.wantError {
				t.Fatalf("CalculateHours() error = %v, want %v", err, test.wantError)
			}
			if err != nil {
				return
			}
			if !equalHoursMetrics(got, test.want) {
				t.Fatalf("CalculateHours() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestHoursMetricsJSONPreservesExactSignedDecimalsAndUnavailableValues(t *testing.T) {
	value := HoursMetrics{EstimatedHours: decimalHours(800), ActualHours: decimalHours(550), Deviation: decimalHours(-250)}
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if got, want := string(encoded), `{"estimated_hours":8.00,"actual_hours":5.50,"deviation":-2.50}`; got != want {
		t.Fatalf("json.Marshal() = %s, want %s", got, want)
	}

	value.EstimatedHours = nil
	value.Deviation = nil
	encoded, err = json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() with unavailable values error = %v", err)
	}
	if got, want := string(encoded), `{"estimated_hours":null,"actual_hours":5.50,"deviation":null}`; got != want {
		t.Fatalf("json.Marshal() with unavailable values = %s, want %s", got, want)
	}
}

func TestDecimalHoursJSONFormatsNegativeFractions(t *testing.T) {
	encoded, err := json.Marshal(DecimalHours(-50))
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if got, want := string(encoded), `-0.50`; got != want {
		t.Fatalf("json.Marshal() = %s, want %s", got, want)
	}
}

func decimalHours(value int64) *DecimalHours {
	result := DecimalHours(value)
	return &result
}

func equalHoursMetrics(a, b HoursMetrics) bool {
	return equalDecimalHours(a.EstimatedHours, b.EstimatedHours) &&
		equalDecimalHours(a.ActualHours, b.ActualHours) &&
		equalDecimalHours(a.Deviation, b.Deviation)
}

func equalDecimalHours(a, b *DecimalHours) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
