package myerror

import "fmt"

type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation failed on %s with value %s. %s", e.Field, e.Value, e.Message)
}

func NewValidationError(f, v, m string) *ValidationError {
	return &ValidationError{f, v, m}
}
