package metrics

import (
	"testing"

	"github.com/Aguscampos04/Carva/internal/domain/hours"
)

func TestCalculateHours(t *testing.T) {
	tests := []struct {
		name      string
		estimates []*hours.Value
		actual    hours.Value
		want      HoursMetrics
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CalculateHours(test.estimates, test.actual); !equalHoursMetrics(got, test.want) {
				t.Fatalf("CalculateHours() = %+v, want %+v", got, test.want)
			}
		})
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
