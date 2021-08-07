package log

import (
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Manager struct{}

var logManager *Manager

func New() *Manager {
	// configure zero log
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	// initiate gocondor logger struct
	logManager = &Manager{}

	return logManager
}

func Resolve() *Manager {
	return logManager
}

func GetHttpLoggingMiddleware() gin.HandlerFunc {
	logger := logger.SetLogger()

	return logger
}
