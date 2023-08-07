package repository

import (
	"ddd_go_example/internal/app/domain/custom_error"
	"ddd_go_example/internal/app/domain/model/command_model"
	"ddd_go_example/internal/app/domain/model/value_object/user"
)

type UserRepository interface {
	Save(cmd command_model.User) (user.UserId, custom_error.RepositoryError)
	Update(cmd command_model.User) custom_error.RepositoryError
	Delete(cmd command_model.User) custom_error.RepositoryError
	FindByIds(ids []user.UserId) ([]command_model.User, custom_error.RepositoryError)
}
