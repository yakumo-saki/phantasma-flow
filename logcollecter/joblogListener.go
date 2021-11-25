package logcollecter

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type LogListenerModule struct {
	procman.ProcmanModuleStruct
}

func (m *LogListenerModule) IsInitialized() bool {
	return m.Initialized
}

func (m *LogListenerModule) Initialize() error {
	// used for procman <-> module communication
	// procman -> PAUSE(prepare for backup) is considered
	m.Name = "LogListener"
	m.Initialized = true
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *LogListenerModule) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return m.Name
}

func (m *LogListenerModule) Start(inCh <-chan string, outCh chan<- string) error {
	m.FromProcmanCh = inCh
	m.ToProcmanCh = outCh
	log := util.GetLoggerWithSource(m.GetName(), "start")

	log.Info().Msgf("Starting %s server.", m.GetName())

	m.ToProcmanCh <- procman.RES_STARTUP_DONE

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

func (sv *LogListenerModule) Shutdown() {
	log := util.GetLoggerWithSource(sv.GetName(), "shutdown")
	log.Debug().Msg("Shutdown initiated")
	sv.RootCancel()
}
