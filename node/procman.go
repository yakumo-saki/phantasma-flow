package node

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type NodeManager struct {
	procman.ProcmanModuleStruct

	// runQueue []objects.JobDefinition
	// nodes    []objects.NodeDefinition
}

// returns this instance is initialized or not.
// When procman.Add, Procman calls Initialize() if not initialized.
func (m *NodeManager) IsInitialized() bool {
	return m.Initialized
}

// initialize this instance.
// Between Initialize and Start, no shutdown is called when error occures.
// so, dont initialize something needs shutdown sequence.
func (m *NodeManager) Initialize() error {
	// used for procman <-> module communication
	// procman -> PAUSE(prepare for backup) is considered
	m.Name = "NodeManager" // if you want to multiple instance, change name here
	m.Initialized = true
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *NodeManager) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return m.Name
}

// lets roll! Do not forget to save procmanCh from parameter.
func (nm *NodeManager) Start(inCh <-chan string, outCh chan<- string) error {
	nm.FromProcmanCh = inCh
	nm.ToProcmanCh = outCh
	log := util.GetLoggerWithSource(nm.GetName(), "main")

	log.Info().Msgf("Starting %s.", nm.GetName())

	nodeCh := messagehub.Subscribe(messagehub.TOPIC_NODE_DEFINITION, nm.GetName())

	nm.ToProcmanCh <- procman.RES_STARTUP_DONE

	// wait for other message from Procman
	for {
		select {
		case v := <-nm.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case node := <-nodeCh:
			log.Info().Msgf("%s", node)
		case <-nm.RootCtx.Done():
			goto shutdown
		}
	}

shutdown:
	// stop all node

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
	nm.RootCancel()
}
