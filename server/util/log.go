package util

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var ZeroLogger *zerolog.Logger

func GetLogger() zerolog.Logger {
	if ZeroLogger == nil {
		w := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05.999",
		}

		zl := zerolog.New(w).With().Timestamp().Caller().Logger()
		ZeroLogger = &zl
	}

	return *ZeroLogger
}

func GetLoggerWithSource(name ...string) zerolog.Logger {
	return GetLogger().With().Str("source", strings.Join(name, "/")).Logger()
}

func Nvl(check string, replaced string) string {
	if check == "" {
		return replaced
	}
	return check
}
