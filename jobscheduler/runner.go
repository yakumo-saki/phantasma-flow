package jobscheduler

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// job runner
// TODO: use context !!!!!!
func (js *JobScheduler) runner() {
	log := util.GetLoggerWithSource(js.GetName(), "runner")
	for {
		benchTime := time.Now()
		js.mutex.Lock()
		now := time.Now()
		nowUnix := now.Unix()

		// scan for all schedules
		for e := js.schedules.Front(); e != nil; e = e.Next() {
			schedule := e.Value.(schedule)
			if schedule.time <= nowUnix {
				// this is temporary
				log.Debug().Msgf("Running jobId=%s runId=%s", schedule.jobId, schedule.runId)
				js.scheduleWithoutLock(schedule.jobId, now)
			}
		}

		js.mutex.Unlock()

		// over time check
		took := time.Since(benchTime).Milliseconds()
		if took > 500 {
			log.Warn().Msgf("Runner took %d ms", took)
		}
		time.Sleep(1 * time.Second)
	}
}
