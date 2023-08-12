package query_model

import "ddd_go_example/internal/app/domain/model/value_object/value_user"

type User struct {
	UserId   value_user.UserId
	UserName value_user.UserName
	Birthday value_user.UserBirthday
}
