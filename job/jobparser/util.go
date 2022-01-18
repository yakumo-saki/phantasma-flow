package jobparser

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func CreateJobLogMsg(jobstep ExecutableJobStep) *objects.JobLogMessage {
	lm := objects.JobLogMessage{}
	lm.Version = jobstep.Version
	lm.JobId = jobstep.JobId
	lm.RunId = jobstep.RunId
	lm.JobNumber = jobstep.JobNumber
	lm.LogDateTime = util.GetDateTimeString()

	return &lm
}
