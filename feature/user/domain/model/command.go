package model

import (
	"time"

	customerror "go_ddd_example/share/custom_error"
)

type UserCommand struct {
	ID       UserID
	Name     UserName
	CreateAt time.Time
	UpdateAt time.Time
}

func CreateUserCommand(name string) (UserCommand, customerror.ModelError) {
	userName, err := ParseUserName(name)
	if err != nil {
		return UserCommand{}, err
	}

	return UserCommand{
		ID:       NewUserID(),
		Name:     userName,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}, nil
}
