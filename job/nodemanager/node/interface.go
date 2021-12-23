package node

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

// execNodeImpl represents node connection (SSH connection or a process)
// This struct is not recycled, single shot execution.
type execNodeImpl interface {
	// initialize node
	// (example: ssh -> connect)
	Initialize(def objects.NodeDefinition, jobStep jobparser.ExecutableJobStep) error

	// Run jobStep (passed on Initialize). when error occured, job step must be fail.
	// After run, do cleanup.
	Run(ctx context.Context)
}
