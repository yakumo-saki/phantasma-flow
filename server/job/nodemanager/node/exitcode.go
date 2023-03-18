package node

import (
	"errors"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

// PHFLOW: definition error
const EC_DEF_ERR = -20000

// *ssh.ExitMissingError
const EC_MISSING = -20001

// *os.PathError
const EC_PATH_ERR = -20002

// the other error
const EC_OTHER_ERR = -29999

func exitCodeFromError(err error) (int, string) {

	var (
		ee *exec.ExitError
		em *ssh.ExitMissingError
		pe *os.PathError
	)

	if err == nil {
		return 0, "no error"
	} else if errors.As(err, &ee) {
		return ee.ExitCode(), "non-zero exit code"
	} else if errors.As(err, &em) {
		return EC_MISSING, "ExitMissingError"
	} else if errors.As(err, &pe) {
		return EC_PATH_ERR, "PathError, no such file or permission denied"
	}

	return EC_OTHER_ERR, "Unknown error"
}
