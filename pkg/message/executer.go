package message

import "github.com/yakumo-saki/phantasma-flow/pkg/objects"

// From repository to all

// reason
const JOB_STEP_START = "JOB-STEP-START" // reason: JOB STEP START
const JOB_STEP_END = "JOB-STEP-END"     // reason: JOB STEP END
const JOB_END = "JOB-END"               // reason: JOB END
const JOB_START = "JOB-START"           // reason: JOB START

type ExecuterMsg struct {
	Subject  string // notification type.
	JobId    string
	RunId    string
	Version  objects.ObjectVersion
	StepName string // JOB_STEP_* only
	Success  bool   // JOB_END or JOB_STEP_END
	Reason   string // why success = true / false
	ExitCode int    // JOB_STEP_END only
}
