package node

import (
	"github.com/jinzhu/copier"
	"github.com/yakumo-saki/phantasma-flow/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Node struct {
	queue []objects.JobDefinition // FIFO
	objects.NodeDefinition
}

func (n *Node) Start(nodeDef objects.NodeDefinition) {
	log := util.GetLogger()
	copier.Copy(nodeDef, n)

	log.Info().Msgf("Node %s started.", n.Name)
}
