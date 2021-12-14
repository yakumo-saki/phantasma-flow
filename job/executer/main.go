package executer

import (
	"container/list"
	"context"
	"fmt"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Executer struct {
	procman.ProcmanModuleStruct

	mutex     sync.Mutex
	jobQueue  map[string][]jobparser.ExecutableJobStep // map[runId] -> []ExecutableJobStep
	nodeQueue map[string]list.List                     // map[nodeId] -> list.List<ExecutableJobStep>
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
	ex.jobQueue = make(map[string][]jobparser.ExecutableJobStep)
	ex.nodeQueue = make(map[string]list.List)
	ex.Initialized = true
	return nil
}

func (ex *Executer) Start(inCh <-chan string, outCh chan<- string) error {
	ex.FromProcmanCh = inCh
	ex.ToProcmanCh = outCh

	log := util.GetLoggerWithSource(ex.GetName(), "main")
	log.Info().Msgf("Starting %s.", ex.GetName())

	jobEndCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, ex.GetName())

	ex.ToProcmanCh <- procman.RES_STARTUP_DONE

	for {
		select {
		case <-ex.RootCtx.Done():
			goto shutdown
		case msg, ok := <-jobEndCh:
			if !ok {
				continue
			}

			exeMsg := msg.Body.(*message.ExecuterMsg)
			fmt.Println(exeMsg)

			switch exeMsg.Subject {
			case message.JOB_END:
				// job complete then delete from queue
			case message.JOB_STEP_END:
				// step_end then store job result.
				// step_end then exec next step or job abort
			default:
				continue
			}

		}
	}

shutdown:
	messagehub.Unsubscribe(messagehub.TOPIC_JOB_REPORT, ex.GetName())
	log.Debug().Msgf("%s stopped.", ex.GetName())

	ex.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (ex *Executer) Shutdown() {
	log := util.GetLogger()
	log.Debug().Msg("Shutdown initiated")
	ex.RootCancel()
}
