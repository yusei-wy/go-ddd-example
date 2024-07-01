package customerror

import "fmt"

type RepositoryError interface {
	Error() string
}

type ImplRepositoryError struct {
	inner error
	msg   string
}

func NewRepositoryError(inner error) RepositoryError {
	return ImplRepositoryError{inner, ""}
}

func NewRepositoryErrorWithMessage(inner error, msg string) RepositoryError {
	return ImplRepositoryError{inner, msg}
}

func (e ImplRepositoryError) Error() string {
	if e.msg != "" {
		return fmt.Errorf("RepositoryError: %s: %w", e.msg, e.inner).Error()
	}

	return fmt.Errorf("RepositoryError: %w", e.inner).Error()
}
