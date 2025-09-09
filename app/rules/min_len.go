package rules

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/goravel/framework/contracts/validation"
)

type MinLen struct {
	min int
}

func (r *MinLen) Signature() string { return "min_len" }

func (r *MinLen) Passes(_ validation.Data, val any, options ...any) bool {
	if len(options) == 0 {
		r.min = 0
		return false
	}

	n, err := strconv.Atoi(fmt.Sprint(options[0]))
	if err != nil || n < 0 {
		r.min = 0
		return false
	}
	r.min = n

	var s string
	switch v := val.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		// Bukan string → gagal
		return false
	}

	return utf8.RuneCountInString(s) >= r.min
}

func (r *MinLen) Message() string {
	return fmt.Sprintf("The :attribute must be at least %d characters.", r.min)
}
