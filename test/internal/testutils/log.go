package testutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/yakumo-saki/phantasma-flow/global"
	"github.com/yakumo-saki/phantasma-flow/util"
)

var testStartTime *time.Time

func SetupTestLogger(t *testing.T) {
	cw := zerolog.ConsoleWriter{
		Out:        zerolog.NewTestWriter(t),
		NoColor:    true,
		TimeFormat: ":",
	}

	zl := zerolog.New(cw).With().Timestamp().Logger()
	util.ZeroLogger = &zl
}

func StartTest() {
	global.DEBUG = true

	t := time.Now()
	testStartTime = &t

	fmt.Printf("\n\nTEST START %s\n\n", t.Format(time.RFC3339))
}

func EndTest() {
	if testStartTime != nil {
		diff := time.Since(*testStartTime)
		fmt.Printf("\n\nEND TEST.   took %s\n\n", diff.String())
		return
	}
	fmt.Printf("\n\nEND TEST\n\n")
}
