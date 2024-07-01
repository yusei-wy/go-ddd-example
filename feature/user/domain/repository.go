package domain

import (
	"go_ddd_example/feature/user/domain/model"
	customerror "go_ddd_example/share/custom_error"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(cmd model.UserCommand) customerror.RepositoryError
	GetUser(userId uuid.UUID) (*model.User, customerror.RepositoryError)
}
