package rules

import (
	"github.com/goravel/framework/contracts/validation"
)

type MaxLen struct {
}

// Signature The name of the rule.
func (receiver *MaxLen) Signature() string {
	return "max_len"
}

// Passes Determine if the validation rule passes.
func (receiver *MaxLen) Passes(data validation.Data, val any, options ...any) bool {
	return true
}

// Message Get the validation error message.
func (receiver *MaxLen) Message() string {
	return ""
}
