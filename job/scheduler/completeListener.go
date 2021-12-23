package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Create new schedule when job ends.
func (js *JobScheduler) jobCompleteListener(ctx context.Context, startWg, shutdownWg *sync.WaitGroup) {
	const NAME = "jobCompleter"
	log := util.GetLoggerWithSource(js.GetName(), NAME)

	jobReportCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, NAME)

	log.Debug().Msgf("%s started.", NAME)
	startWg.Done()
	defer shutdownWg.Done()

	for {

		var exeMsg *message.ExecuterMsg
		select {
		case <-ctx.Done():
			log.Debug().Msgf("%s/%s stopped.", js.GetName(), NAME)
			goto shutdown
		case msg := <-jobReportCh:
			exeMsg = msg.Body.(*message.ExecuterMsg)
		}

		now := time.Now()

		switch exeMsg.Subject {
		case message.JOB_END:
			// end
			js.schedule(exeMsg.JobId, now)
		default:
			continue
		}
	}

shutdown:
	// messagehub.StopListen(NAME)
	log.Debug().Msgf("%s stopped.", NAME)
}
