package node

import (
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func createJobLogMsg(jobStep jobparser.ExecutableJobStep) *logfile.JobLogMessage {
	msg := logfile.JobLogMessage{}
	msg.JobId = jobStep.JobId
	msg.RunId = jobStep.RunId
	msg.Stage = logfile.LM_STAGE_JOB
	msg.Version = jobStep.Version
	msg.LogDateTime = util.GetDateTimeString()
	return &msg

}
