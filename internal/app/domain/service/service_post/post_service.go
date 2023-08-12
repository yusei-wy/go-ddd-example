package service_post

import (
	"ddd_go_example/internal/app/domain/model/query_model"
	"ddd_go_example/internal/app/domain/repository"
)

type PostService struct {
	conn           repository.DBConnection
	userRepository *repository.PostRepository
}

func NewPostService(conn repository.DBConnection, userRepository *repository.PostRepository) *PostService {
	return &PostService{
		conn:           conn,
		userRepository: userRepository,
	}
}

func (s *PostService) Update(conn repository.DBConnection) error {
	// TODO: implement
	return nil
}

func (s *PostService) Delete(conn repository.DBConnection) error {
	// TODO: implement
	return nil
}

func (s *PostService) FindByIds(conn repository.DBConnection) ([]*query_model.Post, error) {
	// TODO: implement
	return []*query_model.Post{}, nil
}
