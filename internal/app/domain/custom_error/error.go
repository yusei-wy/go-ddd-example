package custom_error

type CustomError struct {
	code    int
	message string
}

func (e CustomError) Error() string {
	return e.message
}

type BusinessRuleError = CustomError
type RepositoryError = CustomError
type ServiceError = CustomError
type UseCaseError = CustomError

func NewCustomError(code int, message string) CustomError {
	return CustomError{
		code:    code,
		message: message,
	}
}

func NewBusinessRuleError(code int, message string) BusinessRuleError {
	return NewCustomError(code, message)
}
func NewRepositoryError(code int, message string) RepositoryError {
	return NewCustomError(code, message)
}
func NewServiceError(code int, message string) ServiceError {
	return NewCustomError(code, message)
}
func NewUseCaseError(code int, message string) UseCaseError {
	return NewCustomError(code, message)
}
