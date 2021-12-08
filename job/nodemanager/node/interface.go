package node

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

type execNodeImpl interface {
	// initialize node (example: ssh -> connect)
	Initialize(def objects.NodeDefinition) error
	// Run job. when error occured, job step must be fail.
	Run(ctx context.Context, jobStep jobparser.ExecutableJobStep)
}
