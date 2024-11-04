package env

import (
	"os"

	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

func NewLocal() Environment {
	return Environment{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    log.DEBUG,
	}
}
