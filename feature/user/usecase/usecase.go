package usecase

import (
	"go_ddd_example/feature/user/domain"
	customerror "go_ddd_example/share/custom_error"
)

type UserUseCase interface {
	CreateUser(input CreateUserInput) customerror.UseCaseError
	GetUser(input GetUserInput) (GetUserOutput, customerror.UseCaseError)
}

type UserUseCaseImpl struct {
	repository domain.UserRepository
	service    domain.UserService
}

func NewUserUseCaseImpl(repository domain.UserRepository, service domain.UserService) UserUseCase {
	return &UserUseCaseImpl{repository: repository, service: service}
}
