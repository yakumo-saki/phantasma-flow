package executer

import (
	"container/list"
	"context"

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
	job := queuedJob{}
	job.StepResults = ex.createStepResults(execJobs)
	job.Context, job.Cancel = context.WithCancel(context.Background())
	job.Steps = ex.listToSlice(execJobs)
	ex.jobQueue[jobstep.RunId] = &job
}

// create slice from list.List
func (ex *Executer) listToSlice(execJobs *list.List) []jobparser.ExecutableJobStep {
	slice := []jobparser.ExecutableJobStep{}
	for e := execJobs.Front(); e != nil; e = e.Next() {
		job := e.Value.(jobparser.ExecutableJobStep)
		slice = append(slice, job)
	}
	return slice
}

func (ex *Executer) createStepResults(execJobs *list.List) map[string]*execJobStepResult {
	result := make(map[string]*execJobStepResult)
	for e := execJobs.Front(); e != nil; e = e.Next() {
		job := e.Value.(jobparser.ExecutableJobStep)
		result[job.Name] = &execJobStepResult{}
	}
	return result
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
