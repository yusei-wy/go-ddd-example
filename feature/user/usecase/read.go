package usecase

import (
	"go_ddd_example/feature/user/domain/model"
	customerror "go_ddd_example/share/custom_error"

	"github.com/google/uuid"
)

type GetUserInput struct {
	ID uuid.UUID `json:"id"`
}

type GetUserOutput struct {
	User *model.User
}

func (u *UserUseCaseImpl) GetUser(input GetUserInput) (GetUserOutput, customerror.UseCaseError) {
	user, err := u.service.GetUser(input.ID)
	if err != nil {
		return GetUserOutput{User: nil}, customerror.NewUseCaseError(
			customerror.UseCaseErrorContextNotFound, customerror.NotFoundError("User", input.ID, err))
	}

	return GetUserOutput{User: user}, nil
}
