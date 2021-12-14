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

func (ex *Executer) Initialize() {
	ex.mutex = sync.Mutex{}
	ex.jobQueue = make(map[string][]jobparser.ExecutableJobStep)
	ex.nodeQueue = make(map[string]list.List)
}

func (ex *Executer) Start(ctx context.Context) {
	log := util.GetLoggerWithSource(ex.GetName(), "main")

	jobStepCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, ex.GetName())

	for {
		select {
		case <-ctx.Done():
			goto shutdown
		case msg := <-jobStepCh:
			exeMsg := msg.Body.(message.ExecuterMsg)
			fmt.Println(exeMsg)
			// job complete then delete from queue
			// step_end then store job result.
			// step_end then exec next step or job abort

		}
	}

shutdown:
	log.Debug().Msgf("%s stopped.", ex.GetName())
}

func (ex *Executer) Shutdown() {
	log := util.GetLogger()
	log.Debug().Msg("Shutdown initiated")
	ex.RootCancel()
}
