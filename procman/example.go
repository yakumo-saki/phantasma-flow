package procman

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

type MinimalProcmanModule struct {
	procmanCh    chan string
	shutdownFlag bool
}

func (m *MinimalProcmanModule) Initialize(procmanCh chan string) error {
	// used for procman <-> module communication
	// procman -> PAUSE(prepare for backup) is considered
	m.procmanCh = procmanCh
	return nil
}

func (m *MinimalProcmanModule) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return "MinimalProcmanModule"
}

func (m *MinimalProcmanModule) Start() error {
	log := util.GetLogger()

	log.Info().Msgf("Starting %s server.", m.GetName())
	m.shutdownFlag = false

	go m.loop()

	for {
		select {
		case v := <-m.procmanCh:
			log.Debug().Msgf("Got request %s", v)
		default:
		}

		if m.shutdownFlag {
			break
		}
	}

	log.Info().Msgf("%s Stopped.", m.GetName())
	m.procmanCh <- RES_SHUTDOWN_DONE
	return nil
}

func (m *MinimalProcmanModule) loop() {
	for {
		time.Sleep(1 * time.Second)
		if m.shutdownFlag {
			break
		}
	}

}

func (sv *MinimalProcmanModule) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLogger()
	log.Info().Msg("Shutdown initiated")
	sv.shutdownFlag = true
}
