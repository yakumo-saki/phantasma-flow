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

func TestMetaListener(t *testing.T) {
	startRepository()

	hub := messagehub_impl.MessageHub{}
	messagehub.SetMessageHub(&hub)
	hub.Initialize()
	hub.StartSender()

	processManager := procman.NewProcessManager(make(chan string, 1))
	processManager.AddService(&metalistener.MetaListener{})
	processManager.Start()

	// time.Sleep(500 * time.Millisecond)

	fmt.Println("start")
	start := createExecuterMsg("job1", "job1", message.JOB_START)
	PostAndWait(messagehub.TOPIC_JOB_REPORT, *start)

	time.Sleep(2 * time.Second)

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
	messagehub.Post(messagehub.TOPIC_JOB_REPORT, msg)

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
	fmt.Printf("SET PHFLOW_HOME = %s", home)
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
