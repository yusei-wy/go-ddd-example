package usecase

import customerror "go_ddd_example/share/custom_error"

type CreateUserInput struct {
	Name string `json:"name"`
}

type CreateUserOutput struct{}

type CreateUserUseCase struct{}

func (u *UserUseCaseImpl) CreateUser(input CreateUserInput) customerror.UseCaseError {
	if err := u.service.CreateUser(input.Name); err != nil {
		return customerror.NewUseCaseError(customerror.UseCaseErrorContextUnexpected, err)
	}

	return nil
}
