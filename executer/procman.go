package executer

import (
	"container/list"
	"context"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Executer struct {
	procman.ProcmanModuleStruct

	Job       objects.JobDefinition
	Node      objects.NodeDefinition
	mutex     sync.Mutex
	runnables *list.List
}

func (m *Executer) IsInitialized() bool {
	return m.Initialized
}

func (m *Executer) Initialize() error {
	m.Name = "Executer"
	m.Initialized = true
	m.runnables = list.New()
	m.mutex = sync.Mutex{}
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *Executer) GetName() string {
	return m.Name
}

func (ex *Executer) Start(inCh <-chan string, outCh chan<- string) error {
	log := util.GetLoggerWithSource(ex.GetName(), "main")
	ex.FromProcmanCh = inCh
	ex.ToProcmanCh = outCh

	log.Info().Msgf("Starting %s server.", ex.GetName())

	// subscribe to messagehub
	jobRunCh := messagehub.Listen(messagehub.TOPIC_JOB_DEFINITION, ex.GetName())

	go ex.runner(ex.RootCtx)

	// start ok
	ex.ToProcmanCh <- procman.RES_STARTUP_DONE

	for {
		select {
		case v := <-ex.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case job := <-jobRunCh:
			log.Debug().Msgf("Got JobDefinitionMsg %s", job)

			// TODO schedule -> get real job
			// jobDefMsg := job.Body.(messagehubObjects.JobDefinitionMsg)
			// jobdef := jobDefMsg.JobDefinition
		case <-ex.RootCtx.Done():
			goto shutdown
		}
	}

shutdown:
	log.Info().Msgf("%s Stopped.", ex.GetName())
	ex.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (sv *Executer) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLogger()
	log.Debug().Msg("Shutdown initiated")
	sv.RootCancel()
}
