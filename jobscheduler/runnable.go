package jobscheduler

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// job runner
func (js *JobScheduler) pickRunnable(ctx context.Context) {
	const NAME = "pickRunnable"
	log := util.GetLoggerWithSource(js.GetName(), NAME)
	for {
		benchTime := time.Now()
		js.mutex.Lock()
		now := time.Now()
		nowUnix := now.Unix()

		// scan for all schedules
		for e := js.schedules.Front(); e != nil; e = e.Next() {
			schedule := e.Value.(schedule)
			if schedule.time <= nowUnix {
				log.Debug().Int64("scheduled", schedule.time).Str("jobId", schedule.jobId).
					Str("runId", schedule.runId).Msg("Running")

				schedule.queuedAt = nowUnix
				js.runnables.PushBack(schedule)
				js.schedules.Remove(e)

				// TODO call executer
				messagehub.Post(messagehub.TOPIC_JOB_RUN_REQUEST, struct{}{})
			}
		}

		js.mutex.Unlock()

		// over time check
		took := time.Since(benchTime).Milliseconds()
		if took > 500 {
			log.Warn().Msgf("%s took %d ms", NAME, took)
		}

		select {
		case <-ctx.Done():
			log.Debug().Msgf("%s/%s stopped.", js.GetName(), NAME)
			return
		case <-time.After(1 * time.Second):
			// nothing to do
		}
	}
}
