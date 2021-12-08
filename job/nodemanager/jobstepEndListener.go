package nodemanager

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// cleanup nodeInstance and Restore capacity
func (nm *NodeManager) jobStepEndListener(ctx context.Context) {
	const NAME = "jobStepEndListener"
	log := util.GetLoggerWithSource(nm.GetName(), NAME)
	for {
		select {
		case <-ctx.Done():
			goto shutdown
		}
	}
shutdown:
	log.Debug().Msgf("%s stopped.", NAME)
}
