package logcollecter

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type LogListenerModule struct {
	procman.ProcmanModuleStruct
}

func (m *LogListenerModule) IsInitialized() bool {
	return m.Initialized
}

func (m *LogListenerModule) Initialize(procmanCh chan string) error {
	// used for procman <-> module communication
	// procman -> PAUSE(prepare for backup) is considered
	m.ProcmanCh = procmanCh
	m.Name = "LogListener"
	m.Initialized = true
	return nil
}

func (m *LogListenerModule) GetName() string {
	// Name of module. must be unique.
	// return fix value indicates this module must be singleton.
	// add secondary instance of this module can cause panic by procman.Add
	return m.Name
}

func (m *LogListenerModule) Start() error {
	log := util.GetLogger()

	log.Info().Msgf("Starting %s server.", m.GetName())
	m.ShutdownFlag = false

	for {
		select {
		case v := <-m.ProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		default:
		}

		if m.ShutdownFlag {
			break
		}

		time.Sleep(procman.MAIN_LOOP_WAIT)
	}

	log.Info().Msgf("%s Stopped.", m.GetName())
	m.ProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (sv *LogListenerModule) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLogger()
	log.Info().Msg("Shutdown initiated")
	sv.ShutdownFlag = true
}

func LogListener(conn net.Conn, shutdown <-chan string, stop chan string, logIn <-chan string) {

	defer conn.Close()
	stopFlag := false

	for {
		select {
		case v := <-stop:
			log.Info().Msgf("STOP signal received %s", v)
			stopFlag = true
		case v := <-shutdown:
			log.Info().Msgf("SHUTDOWN signal received %s", v)
			stopFlag = true
		default:
			log.Debug().Msg("Wait for channel")
			message, more := <-logIn
			if more {
				log.Debug().Str("message", message).Msg("msg received")
				_, err := io.Copy(conn, strings.NewReader(message+"\n"))
				if err != nil {
					log.Debug().Err(err).Msg("Send log failed or connection closed")
					stopFlag = true
				}
			} else {
				log.Debug().Msg("msg channel closed")
				break
			}
			log.Debug().Msg("next loop send_data")
		}

		if stopFlag {
			break
		}

	}

	stop <- "STOP"
	log.Info().Msg("send_data stopped")
}
