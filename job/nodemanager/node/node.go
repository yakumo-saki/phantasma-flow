package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
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

func (n *ExecNode) Initialize(def objects.NodeDefinition) error {
	n.nodeDef = def

	var impl execNodeImpl
	if def.NodeType == objects.NODE_LOCAL {
		lo := localExecNode{}
		impl = &lo
	} else {
		msg := fmt.Sprintf("NodeType %s is unknown", def.NodeType)
		panic(msg)
	}

	err := impl.Initialize(def)
	if err == nil {
		n.node = impl
	}
	return err
}

func (n *ExecNode) Run(ctx context.Context, wg *sync.WaitGroup, jobStep jobparser.ExecutableJobStep) {
	n.Running = true
	n.node.Run(ctx, jobStep)
	n.Running = false
	wg.Done()
}
