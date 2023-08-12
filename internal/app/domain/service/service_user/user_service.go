package service_user

import (
	"ddd_go_example/internal/app/domain/model/query_model"
	"ddd_go_example/internal/app/domain/repository"
)

type UserService struct {
	conn           repository.DBConnection
	userRepository *repository.UserRepository
}

func NewUserService(conn repository.DBConnection, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		conn:           conn,
		userRepository: userRepository,
	}
}

func (s *UserService) Update() error {
	// TODO: implement
	return nil
}

func (s *UserService) Delete() error {
	// TODO: implement
	return nil
}

func (s *UserService) FindByIds() ([]*query_model.User, error) {
	// TODO: implement
	return []*query_model.User{}, nil
}
