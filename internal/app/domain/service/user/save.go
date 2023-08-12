package user

type SaveUserInput struct {
	Name     string
	Birthday string
}

type SaveUserOutput struct {
	UserId string
}

func (s *UserService) Save(input SaveUserInput) (SaveUserOutput, error) {
	// TODO: create command

	// TODO: save user

	// TODO: return user id
	return SaveUserOutput{UserId: ""}, nil
}
