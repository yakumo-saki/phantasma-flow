package executer

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Add new job
func (e *Executer) jobRequestHandler(ctx context.Context) {
	NAME := "jobRequest"
	log := util.GetLoggerWithSource(e.GetName(), NAME)
	log.Debug().Msgf("%s/%s started.", e.GetName(), NAME)

	jobReqCh := messagehub.Listen(messagehub.TOPIC_JOB_RUN_REQUEST, NAME)

	for {
		select {
		case <-ctx.Done():
			goto shutdown
		case req := <-jobReqCh:

			req.Body.(message.)

			e.mutex.Lock()

			log.Debug().Msg("critical section")
			// TODO: ... ???

			e.mutex.Unlock()
		}

	}
shutdown:
	log.Debug().Msgf("%s/%s stopped.", e.GetName(), NAME)
}
