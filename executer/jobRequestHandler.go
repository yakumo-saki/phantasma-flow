package executer

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Add new job
func (ex *Executer) jobRequestHandler(ctx context.Context) {
	NAME := "jobRequest"
	log := util.GetLoggerWithSource(ex.GetName(), NAME)
	log.Debug().Msgf("%s/%s started.", ex.GetName(), NAME)

	jobReqCh := messagehub.Listen(messagehub.TOPIC_JOB_RUN_REQUEST, NAME)

	for {
		select {
		case <-ctx.Done():
			goto shutdown
		case reqMsg := <-jobReqCh:

			req := reqMsg.Body.(message.JobRequest)

			job := jobTask{}
			job.Started = false
			job.RunId = req.RunId
			job.JobId = req.JobId
			job.JobDef = objects.JobDefinition{} // XXX get jobdef from repo
			job.JobStepTasks = []*jobStepTask{}
			job.JobStepState = make(map[string]*jobStepStatus)

			step := jobStepTask{}
			step.Name = "step1"
			job.JobStepTasks = append(job.JobStepTasks, &step)
			job.JobStepState[step.Name] = &jobStepStatus{Running: false, Done: false}

			ex.mutex.Lock()

			// TODO send log. job queued
			log.Debug().Str("jobId", req.JobId).Str("runId", req.RunId).Msg("Request received.")
			ex.jobQueue.PushBack(&job)
			ex.mutex.Unlock()
		}

	}
shutdown:
	log.Debug().Msgf("%s/%s stopped.", ex.GetName(), NAME)
}
