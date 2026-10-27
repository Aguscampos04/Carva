// Package hours represents exact hundredths of an hour, independently of Story Points.
package hours

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Value int64

var ErrInvalid = errors.New("hours must be a positive decimal with at most two decimal places")
var decimal = regexp.MustCompile(`^(0|[1-9][0-9]*)(\.[0-9]{1,2})?$`)

func (v *Value) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if !decimal.MatchString(raw) {
		return ErrInvalid
	}
	whole, frac, _ := strings.Cut(raw, ".")
	frac += strings.Repeat("0", 2-len(frac))
	centText := strings.TrimLeft(whole+frac, "0")
	if centText == "" {
		*v = 0
		return nil
	}
	n, err := strconv.ParseInt(centText, 10, 64)
	if err != nil {
		return ErrInvalid
	}
	*v = Value(n)
	return nil
}
func (v Value) MarshalJSON() ([]byte, error) {
	if v < 0 {
		return nil, ErrInvalid
	}
	return []byte(fmt.Sprintf("%d.%02d", v/100, v%100)), nil
}
func Ptr(v Value) *Value { return &v }
