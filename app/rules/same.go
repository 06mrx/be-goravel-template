package rules

import (
	"fmt"

	"github.com/goravel/framework/contracts/validation"
)

type Same struct {
	otherField string // simpan nama field pembanding
}

// Signature The name of the rule.
func (r *Same) Signature() string {
	return "same"
}

// Passes Determine if the validation rule passes.
// options[0] = nama field lain yang harus dicocokkan
func (r *Same) Passes(data validation.Data, val any, options ...any) bool {
	if len(options) == 0 {
		return false // harus ada field pembanding
	}

	field, ok := options[0].(string)
	if !ok || field == "" {
		return false
	}

	r.otherField = field // simpan untuk message

	// Ambil value field lain
	otherVal, exists := data.Get(field)
	if !exists {
		return false
	}

	// Bandingkan
	return val == otherVal
}

// Message Get the validation error message.
func (r *Same) Message() string {
	return fmt.Sprintf("This field must be the same as '%s'.", r.otherField)
}
