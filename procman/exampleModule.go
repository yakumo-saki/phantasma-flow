package procman

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

type MinimalProcmanModule struct {
	procmanCh    chan string
	shutdownFlag bool
	Name         string // Recommended for debug
	initialized  bool
}

// returns this instance is initialized or not.
// When procman.Add, Procman calls Initialize() if not initialized.
func (m *MinimalProcmanModule) IsInitialized() bool {
	return m.initialized
}

func (m *MinimalProcmanModule) Initialize(procmanCh chan string) error {
	// used for procman <-> module communication
	// procman -> PAUSE(prepare for backup) is considered
	m.procmanCh = procmanCh
	m.Name = "MinimalProcmanModule" // if you want to multiple instance, change name here
	m.initialized = true
	return nil
}

func (m *MinimalProcmanModule) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return m.Name
}

func (m *MinimalProcmanModule) Start() error {
	log := util.GetLogger()

	log.Info().Msgf("Starting %s.", m.GetName())
	m.shutdownFlag = false

	go m.loop()

	// wait for other message from Procman
	for {
		select {
		case v := <-m.procmanCh:
			log.Debug().Msgf("Got request %s", v)
		default:
		}

		if m.shutdownFlag {
			break
		}

		time.Sleep(MAIN_LOOP_WAIT) // Do not want to rush this loop
	}

	log.Info().Msgf("%s Stopped.", m.GetName())
	m.procmanCh <- RES_SHUTDOWN_DONE
	return nil
}

func (m *MinimalProcmanModule) loop() {
	for {
		time.Sleep(MAIN_LOOP_WAIT)
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
	log.Debug().Msg("Shutdown initiated")
	sv.shutdownFlag = true
}
