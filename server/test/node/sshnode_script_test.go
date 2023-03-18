package node

import (
	"context"
	"testing"

	"github.com/huandu/go-assert"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/job/nodemanager"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func TestSSHNodeRunScript(t *testing.T) {
	JOBID := "ssh_script_test"
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
	capa := nodeMan.GetCapacity(NODE)
	log.Debug().Msgf("node %s capacity = %v", NODE, capa)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	js := jobparser.ExecutableJobStep{}
	js.ExecType = objects.JOB_EXEC_TYPE_SCRIPT
	js.JobId = JOBID
	js.RunId = RUNID
	js.Node = NODE
	js.UseCapacity = 1
	js.Name = "Step1"
	js.Script = "hostname\nLANG=C date\nls -1\nsleep 1"
	nodeMan.ExecJobStep(ctx, js)

	capa = nodeMan.GetCapacity(NODE)
	log.Debug().Msgf("node %s capacity = %v", NODE, capa)

	nodeMan.DebugWaitForAllJobsDone()

	// capacity before complete job is 1
	assert.Equal(t, 1, capa)

	// capacity restored
	capa = nodeMan.GetCapacity(NODE)
	assert.Equal(t, 2, capa)

	testutils.EndTest()

	pman.Shutdown()

}
