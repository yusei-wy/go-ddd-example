package user

type SaveUserInput struct {
	Name     string
	Birthday string
}

type SaveUserOutput struct {
	UserId string
}

func (s *UserService) Save(input SaveUserInput) (SaveUserOutput, error) {
	// TODO: implement
	return SaveUserOutput{UserId: ""}, nil
}
