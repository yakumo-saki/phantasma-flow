package logfileexporter

import (
	"context"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
	"golang.org/x/sync/syncmap"
)

type LogFileExporter struct {
	procman.ProcmanModuleStruct
	logChannelsWg sync.WaitGroup
	logChannels   syncmap.Map // [string(runId)] <-chan LogMessage
	logCloseFunc  syncmap.Map // [string(runId)] context.CancelFunc
}

func (m *LogFileExporter) IsInitialized() bool {
	return m.Initialized
}

func (m *LogFileExporter) Initialize() error {
	// used for procman <-> module communication
	// procman -> PAUSE(prepare for backup) is considered
	m.Name = "LogListener"
	m.Initialized = true
	m.logChannels = syncmap.Map{}
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *LogFileExporter) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return m.Name
}

func (m *LogFileExporter) Start(inCh <-chan string, outCh chan<- string) error {
	m.FromProcmanCh = inCh
	m.ToProcmanCh = outCh
	log := util.GetLoggerWithSource(m.GetName(), "start")

	m.logChannelsWg = sync.WaitGroup{}
	m.logChannelsWg.Add(1)
	go m.LogListener(m.RootCtx)

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
	m.logChannelsWg.Wait()
	log.Info().Msgf("%s Stopped.", m.GetName())
	m.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (sv *LogFileExporter) Shutdown() {
	log := util.GetLoggerWithSource(sv.GetName(), "shutdown")
	log.Debug().Msg("Shutdown initiated")
	sv.RootCancel()
}
