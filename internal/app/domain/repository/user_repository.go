package repository

import (
	"ddd_go_example/internal/app/domain/custom_error"
	"ddd_go_example/internal/app/domain/model/command_model"
	"ddd_go_example/internal/app/domain/model/value_object/value_user"
)

type UserRepository interface {
	Save(cmd command_model.User) (value_user.UserId, custom_error.RepositoryError)
	Update(cmd command_model.User) custom_error.RepositoryError
	Delete(cmd command_model.User) custom_error.RepositoryError
	FindByIds(ids []value_user.UserId) ([]command_model.User, custom_error.RepositoryError)
}
