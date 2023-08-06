package service

import (
	"ddd_go_example/internal/app/domain/model/query_model"
	"ddd_go_example/internal/app/domain/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) Save(conn repository.DBConnection) (string, error) {
	// TODO: implement
	return "", nil
}

func (s *UserService) Update(conn repository.DBConnection) error {
	// TODO: implement
	return nil
}

func (s *UserService) Delete(conn repository.DBConnection) error {
	// TODO: implement
	return nil
}

func (s *UserService) FindByIds(conn repository.DBConnection) ([]*query_model.User, error) {
	// TODO: implement
	return []*query_model.User{}, nil
}
