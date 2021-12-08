package nodemanager

import (
	"container/list"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/global"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/job/nodemanager/node"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func (nm *NodeManager) GetCapacity(name string) int {
	return 1
}

// Exec job step.
// Logs and Results can be listen via messagehub
func (nm *NodeManager) ExecJobStep(ctx context.Context, step jobparser.ExecutableJobStep) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	// fetch first node from nodePool
	nd := nm.nodePool[step.Node]
	nodeMeta := nd.Front().Value.(nodeMeta)
	nm.HasEnoughCapacity(nodeMeta, step)

	// new Node instance
	nodeInst := nodeInstance{}
	ctx, cancel := context.WithCancel(ctx)
	nodeInst.Cancel = cancel

	execNode := node.ExecNode{}
	execNode.Initialize(nodeMeta.Def)
	go execNode.Run(ctx, step)
	nodeMeta.RunningInstances[step.GetId()] = nodeInst
}

func (nm *NodeManager) HasEnoughCapacity(nodeMeta nodeMeta, step jobparser.ExecutableJobStep) {
	if nodeMeta.Capacity < step.JobStepDefinition.UseCapacity {
		msg := fmt.Sprintf("Insufficient node capacity req: %v, node: %v node_is_deprecated: %v",
			nodeMeta.Capacity, step.JobStepDefinition.UseCapacity, nodeMeta.Deprecated)
		if global.DEBUG {
			panic(msg)
		} else {
			log.Error().Msg(msg + ", continue anyway")
		}
	}
}

// Add or Change Node defintiion
// Before call, must get mutex lock
func (nm *NodeManager) nodeDefToPool(nodeDef objects.NodeDefinition) {
	log := util.GetLoggerWithSource(nm.GetName(), "NodeDefToPool")

	nd := nodeMeta{}
	nd.Def = nodeDef
	nd.Capacity = nodeDef.Capacity
	nd.Deprecated = false
	nd.RunningInstances = map[string]nodeInstance{}

	ls, ok := nm.nodePool[nodeDef.Name]
	if !ok {
		ls = list.New()
		ls.PushBack(nd)
		nm.nodePool[nd.Def.Name] = ls
		log.Debug().Msgf("New node added. %s", nd.Def.Name)
	} else {
		if ls.Len() > 0 {
			for e := ls.Front(); e != nil; e = e.Next() {
				n := e.Value.(nodeMeta)
				n.Deprecated = true
				log.Debug().Msgf("Changed node definition. %s", nd.Def.Name)
			}
		}
		ls.PushBack(nd)
	}
}
