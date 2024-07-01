package domain

import (
	"go_ddd_example/feature/user/domain/model"
	customerror "go_ddd_example/share/custom_error"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(name string) customerror.ServiceError
	GetUser(id uuid.UUID) (*model.User, customerror.ServiceError)
}

type UserServiceImpl struct {
	repository UserRepository
}

func NewUserServiceImpl(repository UserRepository) UserService {
	return &UserServiceImpl{
		repository: repository,
	}
}

func (s *UserServiceImpl) CreateUser(name string) customerror.ServiceError {
	cmd, err := model.CreateUserCommand(name)
	if err != nil {
		return customerror.NewServiceError(customerror.ServiceErrorContextValidation, err)
	}

	if err := s.repository.CreateUser(cmd); err != nil {
		return customerror.NewServiceError(customerror.ServiceErrorContextRepository, err)
	}

	return nil
}

func (s *UserServiceImpl) GetUser(id uuid.UUID) (*model.User, customerror.ServiceError) {
	user, err := s.repository.GetUser(id)
	if err != nil {
		return nil, customerror.NewServiceError(customerror.ServiceErrorContextRepository, err)
	}

	return user, nil
}
