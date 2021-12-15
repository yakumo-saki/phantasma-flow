package executer

import (
	"sync"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type execJobStepResult struct {
	Started bool // job step is started
	Ended   bool // job step is ended, regardless success or not
	Success bool // job step is success
}

func (ex *Executer) resultCollecter(startWg, stopWg *sync.WaitGroup) {
	const NAME = "resultCollecter"
	log := util.GetLoggerWithSource(ex.GetName(), NAME)
	log.Info().Msgf("Starting %s/%s.", ex.GetName(), NAME)

	jobEndCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, ex.GetName())

	defer stopWg.Done()
	startWg.Done()

	for {
		select {
		case <-ex.RootCtx.Done():
			goto shutdown
		case msg, ok := <-jobEndCh:
			if !ok {
				continue
			}

			exeMsg := msg.Body.(*message.ExecuterMsg)

			switch exeMsg.Subject {
			case message.JOB_END:
				// job complete then delete from queue
			case message.JOB_STEP_END:
				log.Debug().Msgf("Got JOB_STEP_END %v", exeMsg)
				// step_end then store job result.
				// step_end then exec next step or job abort
				ex.mutex.Lock()
				jobq := ex.jobQueue[exeMsg.RunId]
				stepResult := jobq.StepResults[exeMsg.StepName]
				stepResult.Ended = true

				// XXX need exit code threshold
				if exeMsg.ExitCode == 0 {
					stepResult.Success = true
				} else {
					stepResult.Success = false
				}

				ex.mutex.Unlock()
			default:
				continue
			}

		}
	}

shutdown:
	messagehub.Unsubscribe(messagehub.TOPIC_JOB_REPORT, ex.GetName())
	log.Debug().Msgf("%s/%s stopped.", ex.GetName(), NAME)
}
