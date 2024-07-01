package env

import "github.com/labstack/gommon/log"

type Environment struct {
	DatabaseUrl string
	LogLevel    log.Lvl
}
