package logcollecter

const LM_STAGE_PRE = "prerun"
const LM_STAGE_JOB = "job"
const LM_STAGE_POST = "postrun"

type LogMessage struct {
	Stage       string `json:"stage"`       // prerun, job, postrun
	JobStep     string `json:"jobStep"`     // Stage=job only
	Source      string `json:"source"`      // log from where
	LogDateTime string `json:"logDateTime"` // RFC3339 yyyy-mm-ddTHH:MM:SS.nnnn+TZ
	Message     string `json:"message"`
}
