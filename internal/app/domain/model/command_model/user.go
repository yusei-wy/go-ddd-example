package command_model

import "ddd_go_example/internal/app/domain/model/value_object/user"

type User struct {
	UserId   user.UserId
	UserName user.UserName
	Birthday user.UserBirthday
}
