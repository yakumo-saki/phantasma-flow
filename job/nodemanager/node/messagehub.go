package node

import (
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func createJobLogMsg(seqNo uint64, jobStep jobparser.ExecutableJobStep) *logfile.JobLogMessage {
	msg := logfile.JobLogMessage{}
	msg.JobId = jobStep.JobId
	msg.RunId = jobStep.RunId
	msg.Stage = logfile.LM_STAGE_JOB
	msg.SeqNo = seqNo
	msg.Version = jobStep.Version
	msg.LogDateTime = util.GetDateTimeString()
	return &msg

}
