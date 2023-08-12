package user

import (
	value_user "ddd_go_example/internal/app/domain/model/value_object/user"
)

type SaveUserInput struct {
	UserId   value_user.UserId
	UserName value_user.UserName
	Birthday value_user.UserBirthday
}

type SaveUserOutput struct {
	UserId string
}

func (s *UserService) Save(input SaveUserInput) (SaveUserOutput, error) {
	// TODO: create command
	// cmd := cmd_user.SaveUserCommand{
	// 	UserId:   input.UserId,
	// 	UserName: input.UserName,
	// 	Birthday: input.Birthday,
	// }

	// TODO: save user

	// TODO: return user id
	return SaveUserOutput{UserId: ""}, nil
}
