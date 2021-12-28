package metalistener_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/metalog"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
)

func TestBasicMetaListener(t *testing.T) {
	jobId := "basic_test"

	hub, pman := testutils.StartBaseModules()
	meta := metalog.GetInstance()

	pman.AddService(10, meta)
	pman.Start()

	hub.StartSender()

	testutils.StartTest()

	runId := time.Now().Format("basic_test_2006-01-02_150405")

	fmt.Println("start")
	{
		start := createExecuterMsg(jobId, runId, message.JOB_START)
		messagehub.Post(messagehub.TOPIC_JOB_REPORT, *start)
	}

	{
		step1start := createExecuterMsg(jobId, runId, message.JOB_STEP_START)
		step1start.StepName = "step1"
		messagehub.Post(messagehub.TOPIC_JOB_REPORT, *step1start)
	}

	{
		step1end := createExecuterMsg(jobId, runId, message.JOB_STEP_END)
		step1end.StepName = "step1"
		step1end.ExitCode = 0
		messagehub.Post(messagehub.TOPIC_JOB_REPORT, *step1end)
	}

	{
		jobend := createExecuterMsg(jobId, runId, message.JOB_END)
		jobend.Reason = "success"
		messagehub.Post(messagehub.TOPIC_JOB_REPORT, *jobend)
	}

	messagehub.WaitForQueueEmpty("")

	testutils.EndTest()

	pman.Shutdown()

}

func PostAndWait(topic string, body interface{}) {
	messagehub.Post(topic, body)
	fmt.Println(messagehub.GetQueueLength())
	messagehub.WaitForQueueEmpty("")
	fmt.Println(messagehub.GetQueueLength())
}

func createExecuterMsg(jobId, runId, reason string) *message.ExecuterMsg {
	msg := message.ExecuterMsg{}
	msg.JobId = jobId
	msg.RunId = runId
	msg.Subject = reason

	return &msg

}
