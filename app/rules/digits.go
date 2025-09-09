package rules

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/validation"
)

type Digits struct {
	expected int
	isNumber bool
}

// Signature The name of the rule.
func (r *Digits) Signature() string {
	return "digits"
}

// Passes Determine if the validation rule passes.
func (r *Digits) Passes(data validation.Data, val any, options ...any) bool {
	if len(options) == 0 {
		r.expected = 0
		return false
	}

	expectedLength, err := strconv.Atoi(fmt.Sprint(options[0]))
	if err != nil {
		r.expected = 0
		return false
	}
	r.expected = expectedLength

	strVal := fmt.Sprint(val)

	// Cek apakah numeric
	if _, err := strconv.Atoi(strVal); err != nil {
		r.isNumber = false
		return false
	}
	r.isNumber = true

	return len(strVal) == expectedLength
}

// Message Get the validation error message.
func (r *Digits) Message() string {
	if !r.isNumber {
		return "The :attribute must be numeric."
	}
	return fmt.Sprintf("The :attribute must have exactly %d digits.", r.expected)
}
