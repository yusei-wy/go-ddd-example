package user

import (
	"net/http"

	"go_ddd_example/feature/user/usecase"

	customerror "go_ddd_example/share/custom_error"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	usecase usecase.UserUseCase
}

func NewUserHandler(userUsecase usecase.UserUseCase) UserHandler {
	return UserHandler{usecase: userUsecase}
}

func (h *UserHandler) CreateUser(ctx echo.Context) error {
	var input usecase.CreateUserInput
	if err := ctx.Bind(&input); err != nil {
		return customerror.NewHandlerErrorWithMessage(customerror.HandlerErrorContextParseError, err, "Invalid input")
	}

	if err := h.usecase.CreateUser(input); err != nil {
		return customerror.NewHandlerError(customerror.HandlerErrorContextUseCase, err)
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (h *UserHandler) GetUser(ctx echo.Context) error {
	userId, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		return customerror.NewHandlerErrorWithMessage(customerror.HandlerErrorContextParseError, err, "Invalid input")
	}
	input := usecase.GetUserInput{
		Id: userId,
	}

	user, err := h.usecase.GetUser(input)
	if err != nil {
		return customerror.NewHandlerError(customerror.HandlerErrorContextUseCase, err)
	}

	return ctx.JSON(http.StatusOK, user)
}
