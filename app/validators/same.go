package validators

// SameRule memvalidasi bahwa field A sama dengan field B
type SameRule struct{}

// Validate implements interface validation.Rule
func (r SameRule) Validate(field string, value any, params []string, data map[string]any) bool {
	if len(params) == 0 {
		return false
	}

	otherField := params[0]
	otherValue, exists := data[otherField]
	if !exists {
		return false
	}

	return value == otherValue
}

// Message mengembalikan pesan error default
func (r SameRule) Message(field string, params []string) string {
	if len(params) > 0 {
		return field + " must be the same as " + params[0]
	}
	return field + " must match the other field"
}
