package nodemanager

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/job/nodemanager/node"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

type nodeMeta struct {
	// nodeMeta is neary equivalent to NodeDefinition
	// It represents remote or local node

	Def              objects.NodeDefinition
	Capacity         int
	Deprecated       bool                    // new definition is arrived and this node is old.
	RunningInstances map[string]nodeInstance // nodeInstance.Id -> nodeInstance. now Running nodes.
}

type nodeInstance struct {
	// nodeInstance is signle execution on node.

	Node        node.ExecNode // Node instance
	UseCapacity int           // Used capacity by this execution
	Id          string        // NodeName+RandomString
	Cancel      context.CancelFunc
}
