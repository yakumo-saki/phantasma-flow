package repository

import (
	"fmt"
)

func parseYamlPanic(typestr, kind, id, err, filepath string) {
	msg := fmt.Sprintf("Not %s yaml. Kind='%s' id='%s' %s path=%s",
		typestr, kind, id, err, filepath)
	panic(msg)
}
