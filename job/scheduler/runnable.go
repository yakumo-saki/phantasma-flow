package scheduler

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/job/executer"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
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

				schedule.reason = SC_TYPE_SCHEDULE
				schedule.runAt = nowUnix
				js.schedules.Remove(e)

				// JobId -> ExecutableJobSteps
				execJobs, err := jobparser.BuildExecutableJob(schedule.jobId, schedule.runId)
				if err != nil {
					// job fail.
					// error to reason
				} else {
					// exec it
					executer.GetInstance().AddToRunQueue(&execJobs)
				}
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

//ExecImmediate Exec job by jobId. returns runId
func (js *JobScheduler) ExecImmediate(jobId string) string {
	const NAME = "ExecImmediate"
	log := util.GetLoggerWithSource(js.GetName(), NAME)

	now := time.Now()
	nowUnix := now.Unix()

	schedule := schedule{}
	schedule.time = nowUnix
	schedule.jobId = jobId
	schedule.runId = js.generateRunId()
	schedule.reason = SC_TYPE_IMMEDIATE
	schedule.scheduledAt = nowUnix
	js.schedules.PushBack(schedule)

	log.Debug().Int64("scheduled", schedule.time).Str("jobId", schedule.jobId).
		Str("runId", schedule.runId).Msg("Scheduled immediate")

	js.mutex.Lock()
	defer js.mutex.Unlock()

	return schedule.runId
}
