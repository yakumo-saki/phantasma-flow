package executer

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// job runner
func (ex *Executer) runner(ctx context.Context) {
	const NAME = "runner"
	log := util.GetLoggerWithSource(ex.GetName(), NAME)
	for {
		benchTime := time.Now()
		ex.mutex.Lock()
		now := time.Now()

		// scan for all schedules
		for e := ex.jobQueue.Front(); e != nil; e = e.Next() {
			// scan for runnable step
			job := e.Value.(*jobTask)

			// Send job start message
			if !job.Started {
				log.Debug().Str("runId", job.RunId).Str("jobId", job.JobId).
					Msgf("Job Started")
				job.Started = true
				ex.notifyJobReport(job.JobId, job.RunId, message.JOB_START)
			}

			for _, step := range job.JobStepTasks {
				state := job.JobStepState[step.Name]
				if !state.Running {
					log.Debug().Str("runId", job.RunId).Str("jobId", job.JobId).
						Msgf("Run (dummy) job step: %s step: %v", job.JobDef.Name, step.Name)
					state.Running = true
					msg := ex.createExecuterMsg(job)
					msg.Subject = message.JOB_STEP_START
					msg.StepName = step.Name
					messagehub.Post(messagehub.TOPIC_JOB_REPORT, *msg)

					msg2 := *msg
					msg2.Subject = message.JOB_STEP_END
					msg2.ExitCode = 0
					messagehub.Post(messagehub.TOPIC_JOB_REPORT, msg2)
				}
				// TODO run, if runnable
			}

			// Complete all jobs?
			if IsJobDone(job.JobStepState) {
				job.EndAt = now.Unix()
				ex.notifyJobReport(job.JobId, job.RunId, message.JOB_END)
				ex.jobQueue.Remove(e)
			}
		}

		// messagehub.Post(messagehub.TOPIC_JOB_REPORT, struct{}{})

		ex.mutex.Unlock()

		// over time check
		took := time.Since(benchTime).Milliseconds()
		if took > 500 {
			log.Warn().Msgf("%s took %d ms", NAME, took)
		}

		select {
		case <-ctx.Done():
			goto shutdown
		case <-time.After(1 * time.Second):
			// nothing to do
		}
	}

shutdown:
	log.Debug().Msgf("%s/%s stopped.", ex.GetName(), NAME)
}

func (ex *Executer) notifyJobReport(jobId, runId, reason string) {
	msg := message.ExecuterMsg{}
	msg.JobId = jobId
	msg.RunId = runId
	msg.Subject = reason
	messagehub.Post(messagehub.TOPIC_JOB_REPORT, msg)
}

func (ex *Executer) createExecuterMsg(task *jobTask) *message.ExecuterMsg {
	msg := message.ExecuterMsg{}
	msg.JobId = task.JobId
	msg.RunId = task.RunId
	msg.ExitCode = -1
	return &msg
}

func IsJobDone(stepStats map[string]*jobStepStatus) bool {
	for _, stat := range stepStats {
		if !stat.Done {
			return false
		}
	}
	return true
}
