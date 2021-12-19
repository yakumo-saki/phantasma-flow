package objects

import (
	"fmt"
)

const LM_STAGE_PRE = "prerun"
const LM_STAGE_JOB = "job"
const LM_STAGE_POST = "postrun"

type JobLogMessage struct {
	JobId       string        `json:"jobId"`
	RunId       string        `json:"runId"`
	Version     ObjectVersion `json:"version"`           // Version of job definition
	Stage       string        `json:"stage"`             // LM_STAGE_* prerun, job, postrun
	JobStep     string        `json:"jobStep,omitempty"` // Stage=job only
	Source      string        `json:"source"`            // log from where. STAGE_JOB => stdout/stderr, others=>modulename
	LogDateTime string        `json:"logDateTime"`       // RFC3339 yyyy-mm-ddTHH:MM:SS.nnnn+TZ
	SeqNo       uint64        `json:"seqNo,omitempty"`   // log sequence number (optional)
	Message     string        `json:"message"`           // log message
}

func (log *JobLogMessage) String() string {
	msg := fmt.Sprintf("JobLogMsg: JobId=%s(%v) RunId=%s Seq=%v Step=%s Stage=%s Source=%s Msg=%s",
		log.JobId, log.Version, log.RunId, log.SeqNo, log.JobStep, log.Stage, log.Source, log.Message)
	return msg
}

type JobLogData struct {
	DateTime string `json:"dateTime" comment:"RFC3339"`
	Message  string `json:"message"`
}
