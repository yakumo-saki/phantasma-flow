package jobscheduler

import (
	"context"
	"time"

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
		now := time.Now()

		var exeMsg messagehubObjects.ExecuterMsg
		select {
		case <-ctx.Done():
			log.Debug().Msgf("%s/%s stopped.", js.GetName(), NAME)
			goto shutdown
		case msg := <-jobReportCh:
			exeMsg = msg.Body.(messagehubObjects.ExecuterMsg)
		}

		switch exeMsg.Reason {
		case messagehubObjects.JOB_END:
			// end
		case messagehubObjects.JOB_START:
			// start
		default:
			continue
		}

		js.mutex.Lock()

		for e := js.runnables.Front(); e != nil; e = e.Next() {
			schedule := e.Value.(schedule)

			switch exeMsg.Reason {
			case messagehubObjects.JOB_START:
				schedule.runAt = now.Unix()
			case messagehubObjects.JOB_END:
				schedule.endAt = now.Unix()
				js.runnables.Remove(e)
				js.scheduleWithoutLock(schedule.jobId, now)
			}
		}

		js.mutex.Unlock()
	}

shutdown:
	// messagehub.StopListen(NAME)
	log.Debug().Msgf("%s stopped.", NAME)

}
