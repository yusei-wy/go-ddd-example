package customerror

import (
	"errors"
	"fmt"
)

type HandlerErrorContext int

const (
	HandlerErrorContextParseError = iota
	HandlerErrorContextUseCase
	HandlerErrorContextInternalService
)

type HandlerError interface {
	Error() string
	Context() HandlerErrorContext
	Inner() error
	Message() string
}

type ImplHandlerError struct {
	context HandlerErrorContext
	inner   error
	message string // optional
}

func NewHandlerError(context HandlerErrorContext, inner error) HandlerError {
	return ImplHandlerError{context, inner, ""}
}

func NewHandlerErrorWithMessage(context HandlerErrorContext, inner error, msg string) HandlerError {
	return ImplHandlerError{context, inner, msg}
}

func ConvertUseCaseErrorToHandlerError(err error) error {
	if err == nil {
		return nil
	}

	var useCaseError UseCaseError
	if !errors.As(err, &useCaseError) {
		return nil
	}

	//nolint:exhaustive
	switch useCaseError.Context() {
	case UseCaseErrorContextInvalidInput:
		return NewHandlerError(HandlerErrorContextParseError, useCaseError)
	default:
		return NewHandlerError(HandlerErrorContextUseCase, useCaseError)
	}
}

func (e ImplHandlerError) Error() string {
	if e.message != "" {
		return fmt.Errorf("HandlerError: %s: %w", e.message, e.inner).Error()
	}

	return fmt.Errorf("HandlerError: %w", e.inner).Error()
}

func (e ImplHandlerError) Context() HandlerErrorContext {
	return e.context
}

func (e ImplHandlerError) Inner() error {
	return e.inner
}

func (e ImplHandlerError) Message() string {
	return e.message
}
