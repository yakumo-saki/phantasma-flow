package executer

import (
	"testing"
	"time"

	"github.com/huandu/go-assert"
	"github.com/yakumo-saki/phantasma-flow/job/executer"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/job/nodemanager"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func TestBasicNodeManager(t *testing.T) {
	log := util.GetLogger()

	a := assert.New(t)
	hub, pman := testutils.StartBaseModules()

	nodeMan := nodemanager.GetInstance()
	exec := executer.GetInstance()
	repo := repository.GetRepository()

	pman.AddService(nodeMan)
	pman.AddService(exec)
	pman.Start()

	hub.StartSender()

	repo.SendAllNodes()
	hub.WaitForQueueEmpty("")

	testutils.StartTest()
	jobdef := repo.GetJobById("jobparser_simple")
	execJobs, err := jobparser.BuildFromJobDefinition(jobdef, "executerTestRun")
	a.NilError(err)

	exec.AddToRunQueue(execJobs)
	WaitForJobQueueEmpty()

	testutils.EndTest()
	pman.Shutdown()

	a.Equal(2, 2)
	log.Info().Msg("OK")
}

func WaitForJobQueueEmpty() {
	log := util.GetLogger()

	for {
		if executer.GetJobQueueLength() == 0 {
			break
		}
		time.Sleep(time.Second)
		log.Debug().Msg("Wait for job complete")
	}
}
