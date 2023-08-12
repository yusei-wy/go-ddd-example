package repository

import (
	"ddd_go_example/internal/app/domain/model/command_model/cmd_user"
	"ddd_go_example/internal/app/domain/model/query_model"
	"ddd_go_example/internal/app/domain/model/value_object/value_user"
)

type UserRepository interface {
	Save(cmd cmd_user.SaveUserCommand) (value_user.UserId, error)
	Update() error
	Delete() error
	FindByIds(ids []value_user.UserId) ([]query_model.User, error)
}
