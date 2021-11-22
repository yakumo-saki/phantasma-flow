package jobscheduler

import (
	"container/list"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/messagehubObjects"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Create from jobdefinition. Filter out not needed for scheduling.
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
	schedules *list.List
	mutex     sync.Mutex
}

func (m *JobScheduler) IsInitialized() bool {
	return m.Initialized
}

func (m *JobScheduler) Initialize() error {
	m.Name = "JobScheduler"
	m.Initialized = true
	m.jobs = make(map[string]job)
	m.schedules = list.New()
	m.mutex = sync.Mutex{}
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
	js.ShutdownFlag = false

	// subscribe to messagehub
	jobDefCh := messagehub.Listen(messagehub.TOPIC_JOB_DEFINITION, js.GetName())

	// start ok
	js.ToProcmanCh <- procman.RES_STARTUP_DONE

	for {
		select {
		case v := <-js.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case job := <-jobDefCh:
			log.Debug().Msgf("Got JobDefinitionMsg %s", job)

			// TODO JOBS and re-schedule
			jobDefMsg := job.Body.(messagehubObjects.JobDefinitionMsg)
			jobdef := jobDefMsg.JobDefinition
			id := js.addJob(jobdef)
			js.schedule(id)
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

func (js *JobScheduler) schedule(jobId string) {
	js.mutex.Lock()
	defer js.mutex.Unlock()
	for e := js.schedules.Front(); e != nil; e = e.Next() {
		j := e.Value.(job)
		if j.id == jobId {
			js.schedules.Remove(e)
		}
	}

	newSchedule := schedule{}
	uuid4, _ := uuid.NewRandom()
	newSchedule.runId = uuid4.String()
	newSchedule.time = 0 // TODO: FIXME
	js.schedules.PushFront(newSchedule)
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
	sv.ShutdownFlag = true
}

func (js *JobScheduler) RequestHandler() {
	log := util.GetLogger()
	log.Debug().Msg("JobScheduler start")
}
