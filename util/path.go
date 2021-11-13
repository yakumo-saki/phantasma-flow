package util

import "path/filepath"

func JoinPath(base string, paths []string) string {
	ret := base
	for _, v := range paths {
		ret = filepath.Join(ret, v)
	}
	return ret
}
