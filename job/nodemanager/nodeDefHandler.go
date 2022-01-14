package nodemanager

import (
	"container/list"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// nodeDefHandler Add or Change Node defintiion
// Before call, must get mutex lock
func (nm *NodeManager) nodeDefHandler(orgNodeDef objects.NodeDefinition) {
	log := util.GetLoggerWithSource(nm.GetName(), "NodeDefToPool")

	nodeDef, err := SetDefaultNodeDefinition(orgNodeDef)
	if err != nil {
		panic(err)
	}

	nd := nodeMeta{}
	nd.Def = *nodeDef
	nd.Capacity = nodeDef.Capacity
	nd.Deprecated = false
	nd.RunningInstances = map[string]nodeInstance{}

	ls, ok := nm.nodePool[nodeDef.Id]
	if !ok {
		// new node ID
		ls = list.New()
		ls.PushBack(&nd)
		nm.nodePool[nd.Def.Id] = ls
		log.Debug().Msgf("New node added. id=%s(%s) Cap:%v", nd.Def.Id, nd.Def.DisplayName, nd.Capacity)
	} else {
		// exist node ID
		if ls.Len() > 0 {
			for e := ls.Front(); e != nil; e = e.Next() {
				n := e.Value.(*nodeMeta)
				n.Deprecated = true
				log.Debug().Msgf("Changed node definition. %s", nd.Def.Id)
			}
		}
		ls.PushBack(nd)
	}
}
