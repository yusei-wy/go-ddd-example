package service_user

import (
	"ddd_go_example/internal/app/domain/model/value_object/value_user"
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
	// cmd := cmd_user.SaveUserCommand{
	// 	UserId:   input.UserId,
	// 	UserName: input.UserName,
	// 	Birthday: input.Birthday,
	// }

	// TODO: save user
	// userId, err := s.UserRepository.Save(cmd)
	// if err != nil {
	// 	return SaveUserOutput{}, err
	// }

	return SaveUserOutput{UserId: ""}, nil
}
