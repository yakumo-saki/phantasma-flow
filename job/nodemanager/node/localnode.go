package node

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync/atomic"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type localExecNode struct {
	nodeDef objects.NodeDefinition
	jobStep jobparser.ExecutableJobStep
	seqNo   uint64
}

func (n *localExecNode) GetName() string {
	return "localExecNode"
}

func (n *localExecNode) Initialize(def objects.NodeDefinition) error {
	n.nodeDef = def

	return nil
}

func (n *localExecNode) Run(ctx context.Context, jobStep jobparser.ExecutableJobStep) {

	log := util.GetLoggerWithSource(n.GetName(), "run").With().
		Str("jobId", jobStep.JobId).Str("runId", jobStep.RunId).
		Str("node", n.nodeDef.Id).Str("step", jobStep.Name).Logger()

	n.jobStep = jobStep

	var err error
	var cmd *exec.Cmd
	switch jobStep.ExecType {
	case objects.JOB_EXEC_TYPE_COMMAND:
		log.Trace().Msgf("Run command %s", jobStep.Command)
		cmd = exec.CommandContext(ctx, "sh", "-c", jobStep.Command)
	case objects.JOB_EXEC_TYPE_SCRIPT:
		// TODO implement script run #25
		log.Trace().Msgf("Run script %s", jobStep.Script)
		cmd = exec.CommandContext(ctx, "sh", "-c", jobStep.Script)
	default:
		panic(fmt.Sprintf("Unknown execType %s on %s/%s",
			jobStep.ExecType, jobStep.JobId, jobStep.Name))
	}
	stderr, err := cmd.StderrPipe()
	if err == nil {
		go n.PipeToLog(ctx, "stderr", stderr)
	} else {
		log.Err(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err == nil {
		go n.PipeToLog(ctx, "stdout", stdout)
	} else {
		log.Err(err)
	}

	err = cmd.Run() // block until process exit
	if err != nil {
		log.Err(err)
	}

	log.Debug().Msgf("%v", cmd)
}

func (n *localExecNode) PipeToLog(ctx context.Context, name string, pipe io.Reader) {
	// log := util.GetLoggerWithSource(n.GetName(), "run", name)

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		seq := atomic.AddUint64(&n.seqNo, 1)
		logmsg := scanner.Text()

		msg := createJobLogMsg(seq, n.jobStep)
		msg.Source = name
		msg.Message = logmsg
		messagehub.Post(messagehub.TOPIC_JOB_LOG, msg)
	}

}
