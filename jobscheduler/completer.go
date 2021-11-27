package jobscheduler

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/messagehubObjects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Subscribe for job completion and remove from runnable queue
func (js *JobScheduler) jobCompleter(ctx context.Context) {
	const NAME = "jobCompleter"
	log := util.GetLoggerWithSource(js.GetName(), NAME)

	jobReportCh := messagehub.Listen(messagehub.TOPIC_JOB_REPORT, NAME)

	log.Debug().Msgf("%s started.", NAME)

	for {
		// now := time.Now()

		var exeMsg messagehubObjects.ExecuterMsg
		select {
		case <-ctx.Done():
			log.Debug().Msgf("%s/%s stopped.", js.GetName(), NAME)
			goto shutdown
		case msg := <-jobReportCh:
			exeMsg = msg.Body.(messagehubObjects.ExecuterMsg)
		}
		if exeMsg.Reason != messagehubObjects.JOB_END {
			continue
		}

		js.mutex.Lock()

		// complete job, by runId. and set next schedule

		// scan for all schedules
		for e := js.runnables.Front(); e != nil; e = e.Next() {
			//schedule := e.Value.(schedule)
			//js.runnables.Remove(e)
			// this is mockup, in real some notify from msghub
			//js.scheduleWithoutLock(schedule.jobId, now)
		}

		js.mutex.Unlock()
	}

shutdown:
	// messagehub.StopListen(NAME)
	log.Debug().Msgf("%s stopped.", NAME)

}
