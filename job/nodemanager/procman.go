package nodemanager

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type NodeManager struct {
	procman.ProcmanModuleStruct
	inShutdown bool // NodeManager in shutdown state
	mutex      sync.Mutex
	wg         sync.WaitGroup
	// map[nodename] list.List<nodeMeta> only first nodemeta is used
	// 2nd or later nodemetas stored when node setting is changed.
	// and then first nodemeta is deprecated.
	nodePool map[string]*list.List
}

func (m *NodeManager) IsInitialized() bool {
	return m.Initialized
}

func (m *NodeManager) Initialize() error {
	m.wg = sync.WaitGroup{}
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	m.nodePool = map[string]*list.List{}
	m.Initialized = true
	m.inShutdown = false
	return nil
}

func (m *NodeManager) GetName() string {
	return "NodeManager"
}

// lets roll! Do not forget to save procmanCh from parameter.
func (nm *NodeManager) Start(inCh <-chan string, outCh chan<- string) error {
	nm.FromProcmanCh = inCh
	nm.ToProcmanCh = outCh
	log := util.GetLoggerWithSource(nm.GetName(), "main")

	log.Info().Msgf("Starting %s.", nm.GetName())

	nodeCh := messagehub.Subscribe(messagehub.TOPIC_NODE_DEFINITION, nm.GetName())
	jobRepoCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, nm.GetName())

	nm.ToProcmanCh <- procman.RES_STARTUP_DONE

	// wait for other message from Procman
	for {
		select {
		case v := <-nm.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case msg := <-nodeCh:
			nodeDefMsg := msg.Body.(message.NodeDefinitionMsg)
			nm.mutex.Lock()
			nm.nodeDefHandler(nodeDefMsg.NodeDefinition)
			nm.mutex.Unlock()
		case msg := <-jobRepoCh:
			exeMsg := msg.Body.(*message.ExecuterMsg)
			if exeMsg.Subject == message.JOB_STEP_END {
				nm.mutex.Lock()
				nm.cleanUpNodePool(exeMsg)
				nm.mutex.Unlock()
			}
		case <-nm.RootCtx.Done():
			goto shutdown
		}
	}

shutdown:
	// stop all node
	log.Info().Msg("Wait for cancel all jobs...")
	nm.waitForAllJobsStopped()

	log.Info().Msgf("%s Stopped.", nm.GetName())
	nm.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (nm *NodeManager) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLoggerWithSource(nm.GetName(), "shutdown")
	log.Debug().Msg("Shutdown initiated")
	nm.inShutdown = true
	nm.RootCancel()
}

// cleanup nodeInstance and Restore capacity
func (nm *NodeManager) waitForAllJobsStopped() {
	const NAME = "waitForAllJobsStopped"
	log := util.GetLoggerWithSource(nm.GetName(), NAME)

	doneCh := make(chan struct{})
	go func() {
		nm.wg.Wait()
		close(doneCh)
	}()

	select {
	case <-time.After(1 * time.Minute):
		log.Warn().Msgf("Cancel running jobs time out")
	case <-doneCh:
		log.Debug().Msgf("Cancel running jobs done")
	}

}
