package errors

import (
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field string, message string) error {
	return &ValidationError{Field: field, Message: message}
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("Invalid field '%s': %s", err.Field, err.Message)
}

func (err *ValidationError) Is(target error) bool {
	t, ok := target.(*ValidationError)
	if ok {
		return t.Field == err.Field && t.Message == err.Message
	} else {
		return false
	}
}
