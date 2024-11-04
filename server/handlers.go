package server

import (
	"net/http"

	"go_ddd_example/feature/user"
	"go_ddd_example/feature/user/usecase"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterHandlers(e *echo.Echo) {
	injector := do.New(user.Package())

	userHandler := user.NewUserHandler(do.MustInvoke[*usecase.UserUseCaseImpl](injector))

	// NOTE: handler は error を返さないと HandlerFunc と型が一致しない
	e.GET("/health", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})
	e.POST("/api/v1/private/users", userHandler.CreateUser)
	e.GET("/api/v1/private/users/:userId", userHandler.GetUser)
}
