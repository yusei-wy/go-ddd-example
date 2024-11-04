package env

import "github.com/labstack/gommon/log"

type Environment struct {
	DatabaseURL string
	LogLevel    log.Lvl
}
