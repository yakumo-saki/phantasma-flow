package node

import (
	"bufio"
	"context"
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
		Str("jobId", jobStep.JobId).Str("runId", jobStep.RunId).Str("step", jobStep.Name).Logger()

	log.Debug().Msgf("Run %s", jobStep)

	var err error
	cmd := exec.CommandContext(ctx, "sh", "-c", jobStep.Command)
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

	err = cmd.Run()
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
		// log.Info().Str("name", name).Uint64("seqNo", seq).Msg(logmsg)

		msg := createJobLogMsg(seq, n.jobStep)
		msg.Source = name
		msg.Message = logmsg
		messagehub.Post(messagehub.TOPIC_JOB_LOG, msg)
	}

}
