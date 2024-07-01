package server

import (
	"net/http"

	"go_ddd_example/feature/user"
	"go_ddd_example/share/usecase"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo, useCase usecase.UseCaseFacade) {
	userHandler := user.NewUserHandler(useCase.UserUseCase)

	// NOTE: handler は error を返さないと HandlerFunc と型が一致しない
	e.GET("/health", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})
	e.POST("/api/v1/private/users", userHandler.CreateUser)
	e.GET("/api/v1/private/users/:userId", userHandler.GetUser)
}
