package executer

import (
	"container/list"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func (ex *Executer) AddToRunQueue(execJobs *list.List) {
	if execJobs.Len() == 0 {
		panic("Execute Job is empty.")
	}

	jobstep := execJobs.Front().Value.(jobparser.ExecutableJobStep)

	// send job start message
	msg := ex.createExecuterMsg(jobstep, message.JOB_START)
	messagehub.Post(messagehub.TOPIC_JOB_REPORT, msg)

	// send job start log
	logmsg := ex.createJobLogMsg(jobstep)
	logmsg.Stage = logfile.LM_STAGE_PRE
	logmsg.Message = "Job started."
	messagehub.Post(messagehub.TOPIC_JOB_LOG, logmsg)

	ex.mutex.Lock()
	defer ex.mutex.Unlock()
	ex.jobQueue[jobstep.RunId] = ex.listToSlice(execJobs)
}

func (ex *Executer) listToSlice(execJobs *list.List) []jobparser.ExecutableJobStep {
	slice := make([]jobparser.ExecutableJobStep, execJobs.Len())
	for e := execJobs.Front(); e != nil; e = e.Next() {
		slice = append(slice, e.Value.(jobparser.ExecutableJobStep))
	}
	return slice
}

func (ex *Executer) createExecuterMsg(jobstep jobparser.ExecutableJobStep, subject string) *message.ExecuterMsg {
	msg := message.ExecuterMsg{}
	msg.Version = jobstep.Version
	msg.JobId = jobstep.JobId
	msg.RunId = jobstep.RunId
	msg.Subject = subject

	return &msg

}

func (ex *Executer) createJobLogMsg(jobstep jobparser.ExecutableJobStep) *logfile.JobLogMessage {
	lm := logfile.JobLogMessage{}
	lm.Source = ex.GetName()
	lm.Version = jobstep.Version
	lm.JobId = jobstep.JobId
	lm.RunId = jobstep.RunId
	lm.LogDateTime = util.GetDateTimeString()

	return &lm
}
