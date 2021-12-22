package util

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// MkdirAll calling os.MkdirAll() and if failed then panic
// @param logger allow nil. when nil we call utli.GetLogger() for logger
func MkdirAll(path string, logger *zerolog.Logger) {
	log := logger
	if log == nil {
		loggr := GetLogger()
		log = &loggr
	}

	if e := os.MkdirAll(path, 0750); e != nil {
		log.Err(e).Str("dir", path).Msgf("mkdir failed")
		panic(fmt.Sprintf("mkdir failed %s %s\n", path, e))
	}
}
