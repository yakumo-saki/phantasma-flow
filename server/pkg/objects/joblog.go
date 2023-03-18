package objects

import (
	"fmt"
)

const LM_STAGE_PRE = "prerun"
const LM_STAGE_JOB = "job"
const LM_STAGE_POST = "postrun"

// JobLog data type.
// NOTE: This is not yaml. output is JSON
type JobLogMessage struct {
	JobId       string        `json:"jobId"`
	RunId       string        `json:"runId"`
	JobNumber   int           `json:"jobNumber"`         // Job Sequence Number (from jobmeta. not suitable for key.)
	Version     ObjectVersion `json:"version"`           // Version of job definition
	Stage       string        `json:"stage"`             // LM_STAGE_* prerun, job, postrun
	JobStep     string        `json:"jobStep,omitempty"` // Stage=job only
	Node        string        `json:"node,omitempty"`    // Stage=job only running node id
	Source      string        `json:"source"`            // log from where. STAGE_JOB => stdout/stderr, others=>modulename
	LogDateTime string        `json:"logDateTime"`       // RFC3339 yyyy-mm-ddTHH:MM:SS.nnnn+TZ
	SeqNo       uint64        `json:"seqNo,omitempty"`   // log sequence number (optional)
	Message     string        `json:"message"`           // log message
}

func (log *JobLogMessage) String() string {
	msg := fmt.Sprintf("JobLogMsg: JobId=%s(%v) RunId=%s Seq=%v Step=%s Stage=%s Source=%s Msg=%s",
		log.JobId, log.Version, log.RunId, log.SeqNo,
		log.JobStep, log.Stage, log.Source, log.Message)
	return msg
}

type JobLogData struct {
	DateTime string `json:"dateTime" comment:"RFC3339"`
	Message  string `json:"message"`
}
