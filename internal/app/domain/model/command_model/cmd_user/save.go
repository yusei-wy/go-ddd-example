package cmd_user

import "ddd_go_example/internal/app/domain/model/value_object/value_user"

type SaveUserCommand struct {
	UserId   value_user.UserId
	UserName value_user.UserName
	Birthday value_user.UserBirthday
}

func NewSaveUserCommand(
	userId value_user.UserId,
	userName value_user.UserName,
	birthday value_user.UserBirthday,
) SaveUserCommand {
	return SaveUserCommand{
		UserId:   userId,
		UserName: userName,
		Birthday: birthday,
	}
}
