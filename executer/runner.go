package executer

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// job runner
func (e *Executer) runner(ctx context.Context) {
	const NAME = "runner"
	log := util.GetLoggerWithSource(e.GetName(), NAME)
	for {
		benchTime := time.Now()
		e.mutex.Lock()
		// now := time.Now()

		// scan for all schedules
		for e := e.runnables.Front(); e != nil; e = e.Next() {
		}

		// messagehub.Post(messagehub.TOPIC_JOB_REPORT, struct{}{})

		e.mutex.Unlock()

		// over time check
		took := time.Since(benchTime).Milliseconds()
		if took > 500 {
			log.Warn().Msgf("%s took %d ms", NAME, took)
		}

		select {
		case <-ctx.Done():
			log.Debug().Msgf("%s/%s stopped.", e.GetName(), NAME)
			return
		case <-time.After(1 * time.Second):
			// nothing to do
		}
	}
}
