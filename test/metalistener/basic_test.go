package metalistener_test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/logcollecter/metalistener"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/repository"
)

func TestBasicMetaListener(t *testing.T) {
	jobId := "basic_test"

	startRepository()

	hub := messagehub_impl.MessageHub{}
	hub.Initialize()
	messagehub.SetMessageHub(&hub)

	processManager := procman.NewProcessManager(make(chan string, 1))
	processManager.AddService(&metalistener.MetaListener{})
	processManager.Start()

	hub.StartSender()

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
		jobend.JobResult = "success"
		messagehub.Post(messagehub.TOPIC_JOB_REPORT, *jobend)
	}

	time.Sleep(time.Hour)

	messagehub.WaitForQueueEmpty("")
	fmt.Println("end")
	processManager.Shutdown()

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
	msg.Reason = reason

	return &msg

}

func startRepository() *repository.Repository {

	_, file, _, _ := runtime.Caller(0)

	dir := path.Dir(file)
	for {
		if path.Base(dir) == "test" {
			break
		}
		dir, _ = path.Split(dir)
	}

	home := path.Join(dir, "phantasma-flow")
	fmt.Printf("SET PHFLOW_HOME = %s\n", home)
	os.Setenv("PHFLOW_HOME", home)

	repo := repository.GetRepository()
	err := repo.Initialize()
	if err != nil {
		log.Error().Err(err).Msg("Error occured at reading initialize data")
		log.Error().Msg("Maybe data is corrupted or misseditted.")
		return nil
	}

	return repo
}
