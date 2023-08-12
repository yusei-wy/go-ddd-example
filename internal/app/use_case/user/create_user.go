package user

import "ddd_go_example/internal/app/domain/service/user"

type CreateUserUseCase struct {
	userService *user.UserService
}

func NewCreateUserUseCase(userService *user.UserService) *CreateUserUseCase {
	return &CreateUserUseCase{userService: userService}
}

func (c *CreateUserUseCase) Execute(input user.SaveUserInput) (user.SaveUserOutput, error) {
	return c.userService.Save(input)
}
