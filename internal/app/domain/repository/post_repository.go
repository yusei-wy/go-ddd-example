package repository

import (
	"ddd_go_example/internal/app/domain/custom_error"
	"ddd_go_example/internal/app/domain/model/command_model"
	"ddd_go_example/internal/app/domain/model/value_object/value_post"
)

type PostRepository interface {
	Save(cmd command_model.Post) (value_post.PostId, custom_error.RepositoryError)
	Update(cmd command_model.Post) custom_error.RepositoryError
	Delete(cmd command_model.Post) custom_error.RepositoryError
	FindByIds(ids []value_post.PostId) ([]command_model.Post, custom_error.RepositoryError)
}
