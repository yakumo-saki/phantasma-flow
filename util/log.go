package util

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func GetLogger() zerolog.Logger {
	//		TimeFormat: "2006-01-02T15:04:05.999Z07:00",
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05.999",
	}
	log := zerolog.New(output).
		With().Timestamp().
		Caller().
		Logger()
	return log
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
