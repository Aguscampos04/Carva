package metrics

import (
	"errors"
	"fmt"

	"github.com/Aguscampos04/Carva/internal/domain/hours"
)

var ErrHoursOverflow = errors.New("hour metric exceeds supported range")

// DecimalHours stores hundredths of an hour and permits signed deviations.
type DecimalHours int64

func (value DecimalHours) MarshalJSON() ([]byte, error) {
	amount := int64(value)
	prefix := ""
	if amount < 0 {
		prefix = "-"
		amount = -(amount + 1)
		amount++
	}
	return []byte(fmt.Sprintf("%s%d.%02d", prefix, uint64(amount)/100, uint64(amount)%100)), nil
}

type HoursMetrics struct {
	EstimatedHours *DecimalHours `json:"estimated_hours"`
	ActualHours    *DecimalHours `json:"actual_hours"`
	Deviation      *DecimalHours `json:"deviation"`
}

// CalculateHours sums known estimates and derives the signed actual-minus-estimated deviation.
// A missing estimate makes the estimate sum and deviation unavailable, while actual remains usable.
func CalculateHours(estimates []*hours.Value, actual hours.Value) (HoursMetrics, error) {
	estimated := int64(0)
	available := true
	for _, estimate := range estimates {
		if estimate == nil {
			available = false
			continue
		}
		value := int64(*estimate)
		if value < 0 || estimated > int64(^uint64(0)>>1)-value {
			return HoursMetrics{}, ErrHoursOverflow
		}
		estimated += value
	}
	if actual < 0 {
		return HoursMetrics{}, ErrHoursOverflow
	}

	result := HoursMetrics{ActualHours: decimalPointer(int64(actual))}
	if available {
		result.EstimatedHours = decimalPointer(estimated)
		result.Deviation = decimalPointer(int64(actual) - estimated)
	}
	return result, nil
}

func decimalPointer(value int64) *DecimalHours {
	result := DecimalHours(value)
	return &result
}
