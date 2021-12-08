package logfile

import "github.com/yakumo-saki/phantasma-flow/pkg/objects"

const LM_STAGE_PRE = "prerun"
const LM_STAGE_JOB = "job"
const LM_STAGE_POST = "postrun"

type JobLogMessage struct {
	JobId       string                `json:"jobId"`
	RunId       string                `json:"runId"`
	Version     objects.ObjectVersion `json:"version"`     // Version of job definition
	Stage       string                `json:"stage"`       // LM_STAGE_* prerun, job, postrun
	JobStep     string                `json:"jobStep"`     // Stage=job only
	Source      string                `json:"source"`      // log from where. STAGE_JOB => stdout/stderr, others=>modulename
	LogDateTime string                `json:"logDateTime"` // RFC3339 yyyy-mm-ddTHH:MM:SS.nnnn+TZ
	Message     string                `json:"message"`     // log message
}

type JobLogData struct {
	DateTime string `json:"dateTime" comment:"RFC3339"`
	Message  string `json:"message"`
}
