package procmanExample

import (
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
	m.Name = "MinimalProcmanModule" // if you want to multiple instance, change name here
	m.Initialized = true
	return nil
}

func (m *MinimalProcmanModule) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return m.Name
}

// lets roll! Do not forget to save procmanCh from parameter.
func (m *MinimalProcmanModule) Start(inCh <-chan string, outCh chan<- string) error {
	m.FromProcmanCh = inCh
	m.ToProcmanCh = outCh
	log := util.GetLogger()

	log.Info().Msgf("Starting %s.", m.GetName())
	m.ShutdownFlag = false

	go m.loop()
	m.ToProcmanCh <- procman.RES_STARTUP_DONE

	// wait for other message from Procman
	for {
		select {
		case v := <-m.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		default:
		}

		if m.ShutdownFlag {
			break
		}

		time.Sleep(procman.MAIN_LOOP_WAIT) // Do not want to rush this loop
	}

	log.Info().Msgf("%s Stopped.", m.GetName())
	m.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (m *MinimalProcmanModule) loop() {
	for {
		time.Sleep(procman.MAIN_LOOP_WAIT)
		if m.ShutdownFlag {
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
	sv.ShutdownFlag = true
}
