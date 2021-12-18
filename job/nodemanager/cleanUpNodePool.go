package nodemanager

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// cleanUpNodePool cleanup nodeInstance and Restore capacity
func (nm *NodeManager) cleanUpNodePool(exeMsg *message.ExecuterMsg) {
	const NAME = "cleanUpNodePool"
	log := util.GetLoggerWithSource(nm.GetName(), NAME)

	pool := nm.nodePool[exeMsg.Node]

	for e := pool.Front(); e != nil; e = e.Next() {
		meta := e.Value.(*nodeMeta)
		for instId, instance := range meta.RunningInstances {
			if instance.Node.Running {
				continue
			}

			instance.Cancel()
			meta.Capacity = meta.Capacity + instance.UseCapacity
			if meta.Capacity > meta.Def.Capacity {
				log.Info().Msgf("Node %s capacity %v is over max capacity %v",
					meta.Def.Id, meta.Capacity, meta.Def.Capacity)
				meta.Capacity = meta.Def.Capacity
			}

			delete(meta.RunningInstances, instId)
		}
		log.Trace().Msgf("Cleanup done. %s node RunningInstances=%v",
			meta.Def.Id, len(meta.RunningInstances))
	}
}
