package user

import "ddd_go_example/internal/app/domain/service/service_user"

type CreateUserUseCase struct {
	userService *service_user.UserService
}

func NewCreateUserUseCase(userService *service_user.UserService) *CreateUserUseCase {
	return &CreateUserUseCase{userService: userService}
}

func (c *CreateUserUseCase) Execute(input service_user.SaveUserInput) (service_user.SaveUserOutput, error) {
	return c.userService.Save(input)
}
