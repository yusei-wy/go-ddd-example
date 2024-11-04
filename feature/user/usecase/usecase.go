package usecase

import (
	"go_ddd_example/feature/user/domain"
	customerror "go_ddd_example/share/custom_error"
)

type UserUseCase interface {
	CreateUser(input CreateUserInput) customerror.UseCaseError
	GetUser(input GetUserInput) (GetUserOutput, customerror.UseCaseError)
}

var _ UserUseCase = (*UserUseCaseImpl)(nil)

type UserUseCaseImpl struct {
	service    domain.UserService
	repository domain.UserRepository
}

func NewUserUseCaseImpl(service domain.UserService, repository domain.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{service: service, repository: repository}
}
