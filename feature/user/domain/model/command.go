package model

import (
	"time"

	customerror "go_ddd_example/share/custom_error"
)

//go:generate getters
type UserCommand struct {
	id       UserId
	name     UserName
	createAt time.Time
	updateAt time.Time
}

func CreateUserCommand(name string) (UserCommand, customerror.ModelError) {
	userName, err := ParseUserName(name)
	if err != nil {
		return UserCommand{}, err
	}

	return UserCommand{
		id:       NewUserId(),
		name:     userName,
		createAt: time.Now(),
		updateAt: time.Now(),
	}, nil
}
