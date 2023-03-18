package out

import (
	"fmt"
	"os"
)

func Msg(param ...interface{}) {
	fmt.Fprintln(os.Stderr, param...)
}

func Msgf(format string, param ...interface{}) {
	msg := fmt.Sprintf(format, param...)
	fmt.Fprint(os.Stderr, msg)
}
