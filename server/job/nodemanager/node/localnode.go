package node

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
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

	seqNo      uint64 // log sequence no (use atomic.Add)
	scriptPath string // created script file path.
}

func (n *localExecNode) GetName() string {
	return "localExecNode"
}

func (n *localExecNode) Initialize(def objects.NodeDefinition, jobStep jobparser.ExecutableJobStep) error {
	n.nodeDef = def
	n.jobStep = jobStep

	// create script. if jobStep is SCRIPT
	if jobStep.ExecType == objects.JOB_EXEC_TYPE_SCRIPT {
		var err error
		n.scriptPath, err = createScriptFile(jobStep)
		if err != nil {
			panic(err) // XXX job fail
		}
	}

	return nil
}

// Run runs initialized jobstep.
//  Returns exitcode (may negative value, example -1 means signal killed)
func (n *localExecNode) Run(ctx context.Context) int {

	jobStep := n.jobStep

	log := util.GetLoggerWithSource(n.GetName(), "run").With().
		Str("jobId", n.jobStep.JobId).Str("runId", jobStep.RunId).
		Str("node", n.nodeDef.Id).Str("step", jobStep.Name).Logger()

	var err error
	var cmd *exec.Cmd
	switch n.jobStep.ExecType {
	case objects.JOB_EXEC_TYPE_COMMAND:
		log.Trace().Msgf("Run command %s", jobStep.Command)
		cmd = exec.CommandContext(ctx, "sh", "-c", jobStep.Command)
	case objects.JOB_EXEC_TYPE_SCRIPT:
		// Run script created on initialize #25
		log.Trace().Msgf("Run script %s", n.scriptPath)
		// TODO TEST without context
		// cmd = exec.CommandContext(ctx, n.scriptPath)
		cmd = exec.Command(n.scriptPath)
	default:
		panic(fmt.Sprintf("Unknown execType %s on %s/%s",
			jobStep.ExecType, jobStep.JobId, jobStep.Name))
	}
	stderr, err := cmd.StderrPipe()
	if err == nil {
		go n.pipeToLog("stderr", stderr)
	} else {
		log.Err(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err == nil {
		go n.pipeToLog("stdout", stdout)
	} else {
		log.Err(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Err(err).Msg("cmd.Start() caused error")
	}

	err = cmd.Wait()
	code, msg := exitCodeFromError(err)
	log.Debug().Err(err).Msgf("Exitcode: %v msg: %s", code, msg)

	if n.scriptPath != "" {
		os.Remove(n.scriptPath)
	}

	return code
}

func (n *localExecNode) pipeToLog(name string, pipe io.Reader) {
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
