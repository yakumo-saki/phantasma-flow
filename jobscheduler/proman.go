package jobscheduler

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func (m *JobScheduler) IsInitialized() bool {
	return m.Initialized
}

func (m *JobScheduler) Initialize() error {
	m.Name = "JobScheduler"
	m.Initialized = true
	m.jobs = make(map[string]job)
	m.schedules = list.New()
	m.runnables = list.New()
	m.mutex = sync.Mutex{}
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *JobScheduler) GetName() string {
	return m.Name
}

func (js *JobScheduler) Start(inCh <-chan string, outCh chan<- string) error {
	log := util.GetLoggerWithSource(js.GetName(), "main")
	js.FromProcmanCh = inCh
	js.ToProcmanCh = outCh

	log.Info().Msgf("Starting %s server.", js.GetName())

	// subscribe to messagehub
	jobDefCh := messagehub.Listen(messagehub.TOPIC_JOB_DEFINITION, js.GetName())

	go js.pickRunnable(js.RootCtx)
	go js.jobCompleter(js.RootCtx)

	// start ok
	js.ToProcmanCh <- procman.RES_STARTUP_DONE

	for {
		select {
		case v := <-js.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case job := <-jobDefCh:
			log.Debug().Msgf("Got JobDefinitionMsg %s", job)

			// TODO JOBS and re-schedule
			jobDefMsg := job.Body.(message.JobDefinitionMsg)
			jobdef := jobDefMsg.JobDefinition
			id := js.addJob(jobdef)
			js.schedule(id, time.Now())
		case <-js.RootCtx.Done():
			goto shutdown
		}
	}

shutdown:
	log.Info().Msgf("%s Stopped.", js.GetName())
	js.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

// Add new job
func (js *JobScheduler) addJob(jobDef objects.JobDefinition) string {
	j := job{}
	j.id = jobDef.Id
	j.jobMeta = jobDef.JobMeta
	j.lastRun = 0
	j.name = jobDef.Name

	js.mutex.Lock()
	defer js.mutex.Unlock()
	js.jobs[j.id] = j
	return j.id
}

func (sv *JobScheduler) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLogger()
	log.Debug().Msg("Shutdown initiated")
	sv.RootCancel()
}

func (js *JobScheduler) RequestHandler() {
	log := util.GetLogger()
	log.Debug().Msg("JobScheduler start")
}
