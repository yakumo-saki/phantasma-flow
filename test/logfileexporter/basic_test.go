package logfileexporter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/logexporter/logfileexporter"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func TestBasicLogFileExporter(t *testing.T) {
	jobId := "logfile"

	hub, pman := testutils.StartBaseModules()

	pman.AddService(10, &logfileexporter.LogFileExporter{})
	pman.Start()

	hub.StartSender()

	runId := time.Now().Format("test")
	fmt.Println("start")

	{
		msg := createJobLogMsg(jobId, runId, logfile.LM_STAGE_PRE)
		msg.Source = "executer"
		msg.Message = "Job queued"
		messagehub.Post(messagehub.TOPIC_JOB_LOG, *msg)
	}
	{
		msg := createJobLogMsg(jobId, runId, logfile.LM_STAGE_JOB)
		msg.Source = "executer"
		msg.Message = "Start step1 on node1"
		messagehub.Post(messagehub.TOPIC_JOB_LOG, *msg)
	}
	{
		msg := createJobLogMsg(jobId, runId, logfile.LM_STAGE_JOB)
		msg.Source = "job"
		msg.Message = "hello world"
		messagehub.Post(messagehub.TOPIC_JOB_LOG, *msg)
	}
	{
		msg := createJobLogMsg(jobId, runId, logfile.LM_STAGE_JOB)
		msg.Source = "executer"
		msg.Message = "End step1, exitcode=0"
		messagehub.Post(messagehub.TOPIC_JOB_LOG, *msg)
	}

	messagehub.WaitForQueueEmpty("")
	fmt.Println("end")
	pman.Shutdown()

}

func PostAndWait(topic string, body interface{}) {
	messagehub.Post(topic, body)
	fmt.Println(messagehub.GetQueueLength())
	messagehub.WaitForQueueEmpty("")
	fmt.Println(messagehub.GetQueueLength())
}

func createJobLogMsg(jobId, runId, stage string) *logfile.JobLogMessage {
	msg := logfile.JobLogMessage{}
	msg.JobId = jobId
	msg.RunId = runId
	msg.Stage = stage
	msg.Version = objects.ObjectVersion{Major: 12, Minor: 89}
	msg.LogDateTime = util.GetDateTimeString()
	return &msg

}
