package rules

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/validation"
)

// DigitsBetween rule: field must be numeric and length between min and max
type DigitsBetween struct{}

func (r *DigitsBetween) Signature() string {
	return "digits_between"
}

func (r *DigitsBetween) Passes(data validation.Data, val any, options ...any) bool {
	if len(options) < 2 {
		return false
	}

	min, err1 := strconv.Atoi(fmt.Sprint(options[0]))
	max, err2 := strconv.Atoi(fmt.Sprint(options[1]))
	if err1 != nil || err2 != nil {
		return false
	}

	strVal := fmt.Sprint(val)
	// Pastikan hanya angka
	if _, err := strconv.Atoi(strVal); err != nil {
		return false
	}

	length := len(strVal)
	return length >= min && length <= max
}

func (r *DigitsBetween) Message() string {
	return "The :attribute must be numeric and have between " + ":min" + " and " + ":max" + " digits."
}
