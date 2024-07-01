package customerror

import (
	"errors"
	"fmt"
)

type UseCaseErrorContext int

const (
	UseCaseErrorContextUnexpected UseCaseErrorContext = iota
	UseCaseErrorContextParseError
	UseCaseErrorContextNotFound
	UseCaseErrorContextConflict
	UsecaseErrorContextDatabase
	UseCaseErrorContextInvalidInput
)

type UseCaseError interface {
	Error() string
	Context() UseCaseErrorContext
}

type ImplUseCaseError struct {
	context UseCaseErrorContext
	inner   error
	msg     string
}

func NewUseCaseError(context UseCaseErrorContext, inner error) UseCaseError {
	return ImplUseCaseError{context, inner, ""}
}

func NewUseCaseErrorWithMessage(context UseCaseErrorContext, inner error, msg string) UseCaseError {
	return ImplUseCaseError{context, inner, msg}
}

func (e ImplUseCaseError) Error() string {
	if e.msg != "" {
		return fmt.Errorf("UseCaseError: %s: %w", e.msg, e.inner).Error()
	}

	return fmt.Errorf("UseCaseError: %w", e.inner).Error()
}

func (e ImplUseCaseError) Context() UseCaseErrorContext {
	return e.context
}

func ConvertServiceToUseCaseError(err error) error {
	if err == nil {
		return nil
	}

	var serviceError ServiceError
	if !errors.As(err, &serviceError) {
		return nil
	}

	//nolint:exhaustive
	switch serviceError.Context() {
	case ServiceErrorContextValidation:
		return NewUseCaseError(UseCaseErrorContextInvalidInput, serviceError)
	case ServiceErrorContextRepository:
		return NewUseCaseError(UsecaseErrorContextDatabase, serviceError)
	default:
		return NewUseCaseError(UseCaseErrorContextUnexpected, serviceError)
	}
}
