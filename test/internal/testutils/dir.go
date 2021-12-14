package testutils

import (
	"path"
	"runtime"
	"strings"
)

func GetTestJobDefDir() string {
	fp := path.Join(GetTestPhHomeDir(), "definitions", "job")
	return fp
}

func GetTestPhHomeDir() string {
	fp := path.Join(GetTestBaseDir(), "phantasma-flow")
	return fp
}

func GetTestBaseDir() string {
	_, file, _, _ := runtime.Caller(0)

	count := 0

	dir := path.Dir(file)
	for {
		if path.Base(dir) == "test" {
			break
		}
		dir = strings.TrimRight(dir, "/") // if path ends with "/" path.Split return itself
		dir, _ = path.Split(dir)

		count++
		if count > 50 {
			panic("cannot find test base dir.")
		}
	}

	return dir
}
