package main

import (
	"os"

	"ddd_go_example/cmd/env"
)

func main() {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "local" {
		env.RunLocal()
	} else if appEnv == "staging" {
		env.RunStaging()
	} else if appEnv == "prod" {
		env.RunProd()
	}
}
