package logfileexporter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yakumo-saki/phantasma-flow/logexporter/logfileexporter"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
)

// exporterがタイムアウトして自動クローズした場合の挙動のテスト
// * Alive = false されているか（＝シャットダウン時に Skipされるか）
// * WaitGroup negative しないか
func TestAutoClose(t *testing.T) {

	hub, pman := testutils.StartBaseModules()

	pman.AddService(10, &logfileexporter.LogFileExporter{})
	pman.Start()

	hub.StartSender()

	jobId := "autoclose"
	runId := time.Now().Format("test")
	fmt.Println("start")

	{
		msg := createJobLogMsg(jobId, runId, objects.LM_STAGE_JOB)
		msg.Source = "job"
		msg.Message = "hello world"
		messagehub.Post(messagehub.TOPIC_JOB_LOG, *msg)
	}

	fmt.Println("wait for auto close occures")
	time.Sleep(16 * time.Second)

	messagehub.WaitForQueueEmpty("")
	fmt.Println("end")
	pman.Shutdown()

}
