package executer

import (
	"container/list"
	"context"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Executer struct {
	procman.ProcmanModuleStruct

	mutex     sync.Mutex
	jobQueue  map[string]*queuedJob // map[runId] -> queuedJob
	nodeQueue map[string]list.List  // map[nodeId] -> list.List<ExecutableJobStep>
}

// queuedJob represents a single job.
// Must call Cancel() when job is ended.
type queuedJob struct {
	Steps       []jobparser.ExecutableJobStep // All job steps
	StepResults map[string]*execJobStepResult // map[stepname] -> Job step results
	Context     context.Context               // job step context
	Cancel      context.CancelFunc            // job step context cancel
}

func (ex *Executer) GetName() string {
	return "Executer"
}

func (ex *Executer) IsInitialized() bool {
	return ex.Initialized
}

func (ex *Executer) Initialize() error {
	ex.mutex = sync.Mutex{}
	ex.RootCtx, ex.RootCancel = context.WithCancel(context.Background())
	ex.jobQueue = make(map[string]*queuedJob)
	ex.nodeQueue = make(map[string]list.List)
	ex.Initialized = true
	return nil
}

func (ex *Executer) Start(inCh <-chan string, outCh chan<- string) error {
	ex.FromProcmanCh = inCh
	ex.ToProcmanCh = outCh

	log := util.GetLoggerWithSource(ex.GetName(), "main")
	log.Info().Msgf("Starting %s.", ex.GetName())

	startWg := sync.WaitGroup{}
	stopWg := sync.WaitGroup{}

	startWg.Add(2)
	stopWg.Add(2)
	go ex.resultCollecter(&startWg, &stopWg)
	go ex.queueExecuter(&startWg, &stopWg)

	startWg.Wait()
	ex.ToProcmanCh <- procman.RES_STARTUP_DONE

	<-ex.RootCtx.Done()

	stopWg.Wait()
	log.Debug().Msgf("%s stopped.", ex.GetName())

	ex.ToProcmanCh <- procman.RES_SHUTDOWN_DONE

	return nil
}

func (ex *Executer) Shutdown() {
	log := util.GetLogger()
	log.Debug().Msg("Shutdown initiated")
	ex.RootCancel()
}
