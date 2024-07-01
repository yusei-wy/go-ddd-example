package infrastructure

import (
	"time"

	"go_ddd_example/feature/user/domain/model"

	"github.com/google/uuid"
)

type UserQuery struct {
	Id       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	CreateAt time.Time `db:"created_at"`
	UpdateAt time.Time `db:"updated_at"`
}

func NewUserQuery(cmd model.UserCommand) UserQuery {
	return UserQuery{
		Id:       cmd.Id().Value(),
		Name:     cmd.Name().String(),
		CreateAt: cmd.CreateAt(),
		UpdateAt: cmd.UpdateAt(),
	}
}

type QueryableUser struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
