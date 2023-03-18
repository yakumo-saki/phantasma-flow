package util

import (
	"errors"
	"os"
	"path/filepath"
)

func JoinPath(base string, paths []string) string {
	ret := base
	for _, v := range paths {
		ret = filepath.Join(ret, v)
	}
	return ret
}

func getStat(path string) os.FileInfo {
	stat, err := os.Stat(path)
	if err == nil {
		return stat
	} else if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return stat
}

func IsFileExist(path string) bool {
	stat := getStat(path)
	if stat == nil {
		return false
	} else if stat.IsDir() {
		return false
	}

	return true
}

func IsDirExist(path string) bool {
	stat := getStat(path)
	if stat == nil {
		return false
	} else if stat.IsDir() {
		return true
	}

	return false
}
