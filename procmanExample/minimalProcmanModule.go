package procmanExample

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type MinimalProcmanModule struct {
	procman.ProcmanModuleStruct
}

// returns this instance is initialized or not.
// When procman.Add, Procman calls Initialize() if not initialized.
func (m *MinimalProcmanModule) IsInitialized() bool {
	return m.Initialized
}

// initialize this instance.
// Between Initialize and Start, no shutdown is called when error occures.
// so, dont initialize something needs shutdown sequence.
func (m *MinimalProcmanModule) Initialize() error {
	// used for procman <-> module communication
	// procman -> PAUSE(prepare for backup) is considered
	m.Initialized = true
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *MinimalProcmanModule) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return "MinimalProcmanModule"
}

// lets roll! Do not forget to save procmanCh from parameter.
func (m *MinimalProcmanModule) Start(inCh <-chan string, outCh chan<- string) error {
	m.FromProcmanCh = inCh
	m.ToProcmanCh = outCh
	log := util.GetLogger()

	log.Info().Msgf("Starting %s.", m.GetName())

	go m.loop(m.RootCtx)
	m.ToProcmanCh <- procman.RES_STARTUP_DONE

	// wait for other message from Procman
	for {
		select {
		case v := <-m.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case <-m.RootCtx.Done():
			goto shutdown
		}
	}

shutdown:
	log.Info().Msgf("%s Stopped.", m.GetName())
	m.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (m *MinimalProcmanModule) loop(ctx context.Context) {
	time.Sleep(procman.MAIN_LOOP_WAIT)
	<-ctx.Done()
	log := util.GetLoggerWithSource(m.GetName())
	log.Debug().Msg("loop exit")
}

func (m *MinimalProcmanModule) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLoggerWithSource(m.GetName())
	log.Debug().Msg("Shutdown initiated")
	m.RootCancel()
}
