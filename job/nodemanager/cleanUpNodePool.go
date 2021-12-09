package nodemanager

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// cleanup nodeInstance and Restore capacity
func (nm *NodeManager) cleanUpNodePool(message.ExecuterMsg) {
	const NAME = "cleanUpNodePool"
	log := util.GetLoggerWithSource(nm.GetName(), NAME)

	for _, v := range nm.nodePool {
		for e := v.Front(); e != nil; e = e.Next() {
			meta := e.Value.(nodeMeta)
			for _, instance := range meta.RunningInstances {
				log.Debug().Msg(instance.Id)
			}
		}
	}
}
