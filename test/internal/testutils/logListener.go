package testutils

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type TestLogListener struct {
	procman.ProcmanModuleStruct
}

func (m *TestLogListener) IsInitialized() bool {
	return m.Initialized
}

func (m *TestLogListener) Initialize() error {
	m.Name = "TestLogListener" // if you want to multiple instance, change name here
	m.Initialized = true
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *TestLogListener) GetName() string {
	return m.Name
}

// lets roll! Do not forget to save procmanCh from parameter.
func (m *TestLogListener) Start(inCh <-chan string, outCh chan<- string) error {
	m.FromProcmanCh = inCh
	m.ToProcmanCh = outCh
	log := util.GetLogger()

	log.Info().Msgf("Starting %s.", m.GetName())
	m.ToProcmanCh <- procman.RES_STARTUP_DONE

	ch := messagehub.Subscribe(messagehub.TOPIC_JOB_LOG, m.GetName())

	// wait for other message from Procman
	for {
		select {
		case v := <-m.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case jobLogMsg := <-ch:
			lm := jobLogMsg.Body.(logfile.JobLogMessage)
			log.Info().Msgf("%s", lm)
		case <-m.RootCtx.Done():
			goto shutdown
		}
	}

shutdown:
	log.Info().Msgf("%s Stopped.", m.GetName())
	m.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (m *TestLogListener) Shutdown() {

	log := util.GetLoggerWithSource(m.GetName())
	log.Debug().Msg("Shutdown initiated")
	m.RootCancel()
}
