package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/job/nodemanager"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func TestBasicNodeManager(t *testing.T) {
	JOBID := "manager_test"
	RUNID := "run123"
	hub, pman := testutils.StartBaseModules()

	nodeMan := nodemanager.GetInstance()

	pman.AddService(nodeMan)
	pman.Start()

	hub.StartSender()

	repository.GetRepository().SendAllNodes()
	hub.WaitForQueueEmpty("")

	testutils.StartTest()

	log := util.GetLogger()

	localCap := nodeMan.GetCapacity("local")
	log.Debug().Msgf("node local capacity = %v", localCap)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	js := jobparser.ExecutableJobStep{}
	js.ExecType = objects.JOB_EXEC_TYPE_COMMAND
	js.JobId = JOBID
	js.RunId = RUNID
	js.Name = "Step1"
	js.Command = "date"
	nodeMan.ExecJobStep(ctx, js)

	testutils.EndTest()
	pman.Shutdown()

	assert.Equal(t, 2, localCap)
}
