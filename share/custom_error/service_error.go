package customerror

import "fmt"

type ServiceErrorContext int

const (
	ServiceErrorContextUnexpected ServiceErrorContext = iota
	ServiceErrorContextRepository
	ServiceErrorContextValidation
)

type ServiceError interface {
	Error() string
	Context() ServiceErrorContext
}

type ImplServiceError struct {
	context ServiceErrorContext
	inner   error
	msg     string
}

func NewServiceError(context ServiceErrorContext, inner error) ServiceError {
	return ImplServiceError{context, inner, ""}
}

func NewServiceErrorWithMessage(context ServiceErrorContext, inner error, msg string) ServiceError {
	return ImplServiceError{context, inner, msg}
}

func (e ImplServiceError) Error() string {
	if e.msg != "" {
		return fmt.Errorf("ServiceError: %s: %w", e.msg, e.inner).Error()
	}

	return fmt.Errorf("ServiceError: %w", e.inner).Error()
}

func (e ImplServiceError) Context() ServiceErrorContext {
	return e.context
}
