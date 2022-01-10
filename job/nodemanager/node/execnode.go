package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

type ExecNode struct {
	nodeDef objects.NodeDefinition
	node    execNodeImpl
	Running bool
}

func (n *ExecNode) GetName() string {
	return "ExecNode"
}

func (n *ExecNode) Initialize(def objects.NodeDefinition, jobStep jobparser.ExecutableJobStep) error {
	n.nodeDef = def

	var impl execNodeImpl
	switch def.NodeType {

	case objects.NODE_TYPE_LOCAL:
		node := localExecNode{}
		impl = &node
	case objects.NODE_TYPE_SSH:
		node := sshExecNode{}
		impl = &node
	default:
		msg := fmt.Sprintf("Error Node %s: nodeType %s is unknown. %v", def.Id, def.NodeType, def)
		panic(msg)
	}

	err := impl.Initialize(def, jobStep)
	if err == nil {
		n.node = impl
	}
	return err
}

// Run jobStep context to cancel.
func (n *ExecNode) Run(ctx context.Context, wg *sync.WaitGroup, jobStep jobparser.ExecutableJobStep) {
	n.Running = true
	n.sendJobStepStartMsg(jobStep)

	n.node.Run(ctx)

	n.Running = false
	n.sendJobStepEndMsg(jobStep)
	wg.Done()
}

func (n *ExecNode) sendJobStepStartMsg(jobstep jobparser.ExecutableJobStep) {
	msg := n.createExecuterMsg(jobstep, message.JOB_STEP_START)

	messagehub.Post(messagehub.TOPIC_JOB_REPORT, msg)
}

func (n *ExecNode) sendJobStepEndMsg(jobstep jobparser.ExecutableJobStep) {
	msg := n.createExecuterMsg(jobstep, message.JOB_STEP_END)

	messagehub.Post(messagehub.TOPIC_JOB_REPORT, msg)
}

func (n *ExecNode) createExecuterMsg(jobstep jobparser.ExecutableJobStep, subject string) *message.ExecuterMsg {
	msg := message.ExecuterMsg{}
	msg.Version = jobstep.Version
	msg.JobId = jobstep.JobId
	msg.RunId = jobstep.RunId
	msg.StepName = jobstep.Name
	msg.Node = jobstep.Node
	msg.Subject = subject

	// fmt.Printf("Job REPORT: %s, Job:%s/%s RunId:%s\n", subject, msg.JobId, msg.StepName, msg.RunId)

	return &msg
}
