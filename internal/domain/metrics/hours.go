package metrics

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Aguscampos04/Carva/internal/domain/hours"
)

var ErrHoursOverflow = errors.New("hour metric exceeds supported range")

// DecimalHours stores hundredths of an hour and permits signed deviations.
type DecimalHours int64

var decimalHourPattern = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]{1,2})?$`)

func (value *DecimalHours) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if !decimalHourPattern.MatchString(raw) {
		return ErrHoursOverflow
	}
	negative := strings.HasPrefix(raw, "-")
	if negative {
		raw = strings.TrimPrefix(raw, "-")
	}
	whole, fraction, _ := strings.Cut(raw, ".")
	fraction += strings.Repeat("0", 2-len(fraction))
	wholeValue, err := strconv.ParseUint(whole, 10, 64)
	if err != nil || wholeValue > (^uint64(0)>>1)/100 {
		return ErrHoursOverflow
	}
	minor, _ := strconv.ParseUint(fraction, 10, 8)
	magnitude := wholeValue*100 + minor
	maxInt64 := ^uint64(0) >> 1
	if (!negative && magnitude > maxInt64) || (negative && magnitude > maxInt64+1) {
		return ErrHoursOverflow
	}
	if negative {
		if magnitude == maxInt64+1 {
			*value = DecimalHours(-1 << 63)
		} else {
			*value = DecimalHours(-int64(magnitude))
		}
	} else {
		*value = DecimalHours(magnitude)
	}
	return nil
}

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

type SprintHoursData struct {
	ID          int64
	Goal        string
	Closed      bool
	Estimates   []*hours.Value
	ActualHours hours.Value
}

type ProjectHoursData struct {
	Estimates   []*hours.Value
	ActualHours hours.Value
	Sprints     []SprintHoursData
}

type SprintHoursMetrics struct {
	ID      int64        `json:"id"`
	Goal    string       `json:"goal"`
	Closed  bool         `json:"closed"`
	Metrics HoursMetrics `json:"metrics"`
}

type ProjectHoursMetrics struct {
	Metrics HoursMetrics         `json:"metrics"`
	Sprints []SprintHoursMetrics `json:"sprints"`
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
