package user

import "ddd_go_example/internal/app/domain/model/value_object/user"

type SaveUserCommand struct {
	UserId   user.UserId
	UserName user.UserName
	Birthday user.UserBirthday
}

func NewSaveUserCommand(
	userId user.UserId,
	userName user.UserName,
	birthday user.UserBirthday,
) SaveUserCommand {
	return SaveUserCommand{
		UserId:   userId,
		UserName: userName,
		Birthday: birthday,
	}
}
