// apps/backend/pkg/logger/logger.go

package logger

import (
	"os"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// برای توسعه، از لاگ خوانا استفاده می‌کنیم
	if os.Getenv("GIN_MODE") == "debug" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}