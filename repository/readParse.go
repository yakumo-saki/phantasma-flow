package repository

import (
	"fmt"
)

func parseYamlPanic(typestr, kind, id, err string) {
	msg := fmt.Sprintf("Not %s yaml. Kind='%s' id='%s' %s", typestr, kind, id, err)
	panic(msg)
}
