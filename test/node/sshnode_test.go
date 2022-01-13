package node

import (
	"context"
	"testing"
	"time"

	"github.com/huandu/go-assert"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/job/nodemanager"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func TestSSHNode(t *testing.T) {
	JOBID := "manager_test"
	RUNID := "run123"
	NODE := "sshkeyfile"
	hub, pman := testutils.StartBaseModules()

	nodeMan := nodemanager.GetInstance()

	pman.AddService(80, nodeMan)
	pman.Start()

	hub.StartSender()

	repository.GetRepository().SendAllNodes()
	hub.WaitForQueueEmpty("")

	testutils.StartTest()

	log := util.GetLogger()

	localCap := nodeMan.GetCapacity(NODE)
	log.Debug().Msgf("node local capacity = %v", localCap)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	js := jobparser.ExecutableJobStep{}
	js.ExecType = objects.JOB_EXEC_TYPE_COMMAND
	js.JobId = JOBID
	js.RunId = RUNID
	js.Node = NODE
	js.Name = "Step1"
	js.Command = "echo \"`hostname` today is `date`\""
	nodeMan.ExecJobStep(ctx, js)

	time.Sleep(5 * time.Second)

	testutils.EndTest()

	pman.Shutdown()

	assert.Equal(t, 2, localCap)
}
