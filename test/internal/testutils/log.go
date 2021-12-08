package testutils

import (
	"fmt"
	"time"

	"github.com/yakumo-saki/phantasma-flow/global"
)

var testStartTime *time.Time

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
