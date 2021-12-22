package testutils

import (
	"os"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// PrepareEmptyDir make path empty and exist.
// Any files in path is deleted. If path is not exist, create it
func PrepareEmptyDir(path string) {
	util.MkdirAll(path, nil)
	err := os.RemoveAll(path)
	ErrPanic(err)
	util.MkdirAll(path, nil)
}
