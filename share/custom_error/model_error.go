package customerror

import "fmt"

type ModelError interface {
	Error() string
}

type ImplModelError struct {
	inner error
	msg   string
}

func NewModelError(inner error, msg string) ModelError {
	return ImplModelError{inner, ""}
}

func NewModelErrorWithMessage(inner error, msg string) ModelError {
	return ImplModelError{inner, msg}
}

func (e ImplModelError) Error() string {
	if e.msg != "" {
		return fmt.Errorf("ModelError: %s: %w", e.msg, e.inner).Error()
	}

	return fmt.Errorf("ModelError: %w", e.inner).Error()
}
