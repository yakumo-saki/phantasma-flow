package node

import (
	"container/list"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

type nodePool struct {
	Nodes map[string]list.List // list of node
}

type node struct {
	Def        objects.NodeDefinition
	Capacity   int
	Deprecated bool // new definition is arrived and this node is old.
}
