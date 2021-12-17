package nodemanager

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// cleanUpNodePool cleanup nodeInstance and Restore capacity
func (nm *NodeManager) cleanUpNodePool(exeMsg *message.ExecuterMsg) {
	const NAME = "cleanUpNodePool"
	log := util.GetLoggerWithSource(nm.GetName(), NAME)

	for _, v := range nm.nodePool {
		for e := v.Front(); e != nil; e = e.Next() {
			meta := e.Value.(nodeMeta)
			for _, instance := range meta.RunningInstances {
				// XXX need these
				// instance.Cancel()
				// meta.Capacity = restore
				log.Debug().Msg(instance.Id)
			}
		}
	}
}
