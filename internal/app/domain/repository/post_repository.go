package repository

import (
	"ddd_go_example/internal/app/domain/custom_error"
	cmd_post "ddd_go_example/internal/app/domain/model/command_model/cmd_post"
	query_post "ddd_go_example/internal/app/domain/model/query_model"
	value_post "ddd_go_example/internal/app/domain/model/value_object/value_post"
)

type PostRepository interface {
	Save(cmd cmd_post.SavePostCommand) (value_post.PostId, custom_error.RepositoryError)
	Update() custom_error.RepositoryError
	Delete() custom_error.RepositoryError
	FindByIds(ids []value_post.PostId) ([]query_post.Post, custom_error.RepositoryError)
}
