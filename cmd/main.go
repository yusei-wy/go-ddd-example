package main

import (
	"os"

	"go_ddd_example/app"
	"go_ddd_example/app/env"
)

func main() {
	stage := env.NewStage(os.Getenv("STAGE"))
	app.NewApp(stage)
}
