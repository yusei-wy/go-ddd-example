package repository

import (
	"ddd_go_example/internal/app/domain/custom_error"
	cmd_user "ddd_go_example/internal/app/domain/model/command_model/user"
	"ddd_go_example/internal/app/domain/model/query_model"
	value_user "ddd_go_example/internal/app/domain/model/value_object/user"
)

type UserRepository interface {
	Save(cmd cmd_user.SaveUserCommand) (value_user.UserId, custom_error.RepositoryError)
	Update() custom_error.RepositoryError
	Delete() custom_error.RepositoryError
	FindByIds(ids []value_user.UserId) ([]query_model.User, custom_error.RepositoryError)
}
