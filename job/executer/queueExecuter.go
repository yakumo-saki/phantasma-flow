package executer

import (
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/job/nodemanager"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func (ex *Executer) queueExecuter(startWg, stopWg *sync.WaitGroup) {
	const NAME = "queueExecuter"
	log := util.GetLoggerWithSource(ex.GetName(), NAME)
	log.Info().Msgf("Starting %s/%s.", ex.GetName(), NAME)

	defer stopWg.Done()
	startWg.Done()

	for {
		select {
		case <-ex.RootCtx.Done():
			// XXX need for wait all jobs in running state
			// TODO cancel all step context
			goto shutdown
		case <-time.After(1 * time.Second):
			ex.mutex.Lock()
			for runId, queuedJob := range ex.jobQueue {
				ex.executeRunnable(runId, queuedJob)
			}
			ex.mutex.Unlock()
		}
	}

shutdown:
	log.Debug().Msgf("%s/%s stopped.", ex.GetName(), NAME)
}

func (ex *Executer) executeRunnable(runId string, job *queuedJob) {
	log := util.GetLoggerWithSource(ex.GetName(), "executeRunnable").
		With().Str("runId", runId).Logger()
	for _, step := range job.Steps {
		stat := job.StepResults[step.Name]

		if stat.Started && !stat.Ended {
			// still running. nothing to do
			goto next
		}

		// not started and no presteps (= entrypoint)
		if len(step.PreSteps) == 0 {
			goto runIt
		}

		// check for all PreSteps are done and successful
		for _, pre := range step.PreSteps {
			s, ok := job.StepResults[pre]
			if !ok {
				goto next // no result = not started. not run. should not occur
			}
			if !s.Success {
				goto next // preStep is failed. not run (Job is failed.)
			} else {
				goto runIt
			}
		}

	runIt:
		log.Debug().Msgf("Jobstep start %s/%s", step.JobId, step.Name)
		nodemanager.GetInstance().ExecJobStep(job.Context, step)
		stat.Started = true

	next:
	}

}
