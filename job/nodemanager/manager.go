package nodemanager

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/global"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/job/nodemanager/node"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func (nm *NodeManager) GetCapacity(name string) int {
	const ERRVAL = -1
	log := util.GetLoggerWithSource(nm.GetName(), "GetCapacity")

	if nm.inShutdown {
		log.Warn().Msgf("NodeManager is in shutdown state. No jobs acceptable %s", name)
		return ERRVAL
	}

	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	pool, ok := nm.nodePool[name]

	if !ok {
		log.Warn().Msgf("No node registered for %s", name)
		return ERRVAL
	} else if pool.Len() == 0 {
		log.Warn().Msgf("No node is available for %s", name)
		return ERRVAL
	}

	meta := pool.Front().Value.(nodeMeta)

	return meta.Capacity
}

// Exec job step.
// Logs and Results can be listen via messagehub
func (nm *NodeManager) ExecJobStep(ctx context.Context, step jobparser.ExecutableJobStep) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	// fetch first node from nodePool
	nd, ok := nm.nodePool[step.Node]
	if !ok {
		// XXX JOB FAIL
		log.Error().Msgf("Node '%s' is not found in NodeManager.", step.Node)
		return
	}
	nodeMeta := nd.Front().Value.(nodeMeta)
	nm.HasEnoughCapacity(nodeMeta, step) // check but not stop.

	// new Node instance
	nodeInst := nodeInstance{}
	ctx, cancel := context.WithCancel(ctx)
	nodeInst.Cancel = cancel

	nm.wg.Add(1)

	execNode := node.ExecNode{}
	execNode.Initialize(nodeMeta.Def)
	go execNode.Run(ctx, &nm.wg, step)
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
