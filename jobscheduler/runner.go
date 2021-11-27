package jobscheduler

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// job runner
func (js *JobScheduler) runner(ctx context.Context) {
	const NAME = "runner"
	log := util.GetLoggerWithSource(js.GetName(), NAME)
	for {
		benchTime := time.Now()
		js.mutex.Lock()
		now := time.Now()

		// scan for all schedules
		for e := js.runnables.Front(); e != nil; e = e.Next() {
			schedule := e.Value.(schedule)
			// this is temporary
			// check Node queue space and run
			log.Debug().Int64("scheduled", schedule.time).Str("jobId", schedule.jobId).
				Str("runId", schedule.runId).Msg("Running (and done)")

			schedule.runAt = now.Unix()
			schedule.endAt = now.Unix()
			js.runnables.Remove(e)
			// this is mockup, in real some notify from msghub
			js.scheduleWithoutLock(schedule.jobId, now)
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
