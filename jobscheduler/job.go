package jobscheduler

import (
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type job struct {
	id      string
	name    string
	lastRun int64
	jobMeta objects.JobMetaInfo
}

type schedule struct {
	time  int64  // unixtime
	runId string // sha1 of uuid

}

type JobScheduler struct {
	procman.ProcmanModuleStruct

	jobs      map[string]job
	schedules []job
	mutex     sync.Mutex
}

func (m *JobScheduler) IsInitialized() bool {
	return m.Initialized
}

func (m *JobScheduler) Initialize() error {
	m.Name = "JobScheduler"
	m.Initialized = true
	m.jobs = make(map[string]job)
	m.schedules = make([]job, 50)
	m.mutex = sync.Mutex{}
	return nil
}

func (m *JobScheduler) GetName() string {
	return m.Name
}

func (js *JobScheduler) Start(inCh <-chan string, outCh chan<- string) error {
	log := util.GetLoggerWithSource(js.GetName(), "start")
	js.FromProcmanCh = inCh
	js.ToProcmanCh = outCh

	log.Info().Msgf("Starting %s server.", js.GetName())
	js.ShutdownFlag = false

	// subscribe to messagehub
	msgCh := messagehub.Listen(messagehub.TOPIC_JOB_DEFINITION, js.GetName())

	// start ok
	js.ToProcmanCh <- procman.RES_STARTUP_DONE

	for {
		select {
		case v := <-js.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case job := <-msgCh:
			log.Debug().Msgf("Got request %s", job)
			// TODO JOBS and re-schedule
		default:
		}

		// todo Job Submitting
		if js.ShutdownFlag {
			break
		}

		time.Sleep(procman.MAIN_LOOP_WAIT)
	}

	log.Info().Msgf("%s Stopped.", js.GetName())
	js.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (sv *JobScheduler) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLogger()
	log.Debug().Msg("Shutdown initiated")
	sv.ShutdownFlag = true
}

func (js *JobScheduler) RequestHandler() {
	log := util.GetLogger()
	log.Debug().Msg("JobScheduler start")
}
