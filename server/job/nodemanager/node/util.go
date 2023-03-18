package node

import (
	"fmt"
	"os"
	"strings"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
)

// create script file on phflow node.
//  You need to delete temporary script file after use.
//  return 1: path_to_script 2: error
func createScriptFile(jobStep jobparser.ExecutableJobStep) (string, error) {
	tempFilename := fmt.Sprintf("phflow_temp_%s_%s_*", jobStep.JobId, jobStep.Name)
	tempfile, err := os.CreateTemp("", tempFilename)
	if err != nil {
		return "", err
	}

	// if script has not shebang, /bin/bash assumed
	if !strings.HasPrefix(jobStep.Script, "#!") {
		tempfile.WriteString("#!/bin/bash\n") // XXX #50
	}
	_, err = tempfile.WriteString(jobStep.Script)
	if err != nil {
		return "", err
	}
	err = tempfile.Chmod(os.FileMode(int(0700)))
	if err != nil {
		return "", err
	}

	tempfile.Close()
	if err != nil {
		return "", err
	}

	return tempfile.Name(), nil

}
