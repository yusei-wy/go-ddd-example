package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go_ddd_example/app/env"
	"go_ddd_example/server"
	"go_ddd_example/share/usecase"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApp(stage env.Stage) {
	// switch env {
	// case Debug:
	// 	NewDebug()
	// case Local:
	environment := env.NewLocal()
	// case Staging:
	// 	NewStaging()
	// case Production:
	// 	NewProduction()
	// }

	db, err := sqlx.Open("postgres", environment.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	useCase := usecase.NewUseCaseFacade(db)

	e := echo.New()

	// logger
	e.Logger.SetLevel(environment.LogLevel)

	e.Pre(middleware.RemoveTrailingSlash()) // 末尾の / を削除して URL を統一

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// TODO: CORS
	// TODO: Secure
	// TODO: Static

	// カスタムエラーハンドラ
	e.HTTPErrorHandler = server.CustomHTTPErrorHandler

	server.RegisterHandlers(e, useCase)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start("localhost:8080"); err != nil && !errors.Is(http.ErrServerClosed, err) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Waiting for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
