package nodemanager

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/job/nodemanager/node"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

// nodeMeta is neary equivalent to NodeDefinition
type nodeMeta struct {
	Def              objects.NodeDefinition
	Capacity         int
	Deprecated       bool                    // new definition is arrived and this node is old.
	RunningInstances map[string]nodeInstance // nodeInstance.Id -> nodeInstance. now Running nodes.
}

// nodeInstance
type nodeInstance struct {
	Node   node.ExecNode // Node instance
	Id     string        // NodeName+RandomString
	Cancel context.CancelFunc
}
