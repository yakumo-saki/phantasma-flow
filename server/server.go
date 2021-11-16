package server

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Server struct {
	procmanCh    chan string
	shutdownFlag bool
	listener     net.Listener
}

func (sv *Server) Initialize(procmanCh chan string) error {
	sv.procmanCh = procmanCh
	return nil
}

func (sv *Server) GetName() string {
	return "server"
}

func (sv *Server) startListen() error {
	// Finally start listening
	// TODO change port by config
	psock, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Error().Err(err).Msg("Failed listen tcp.")
		return err
	}
	sv.listener = psock
	return nil
}

func (sv *Server) Start() error {
	log := util.GetLogger()

	log.Info().Msg("Starting socket server.")
	sv.shutdownFlag = false

	err := sv.startListen()
	if err != nil {
		return err
	}
	log.Debug().Msg("TCP Socket start")

	go sv.awaitListener()
	log.Info().Msg("Socket server started.")

	for {
		select {
		case v := <-sv.procmanCh:
			log.Debug().Msgf("Got request from procman %s", v)
		default:
		}

		if sv.shutdownFlag {
			break
		}

		time.Sleep(1 * time.Second)
	}

	sv.listener.Close()
	log.Debug().Msg("Main thread exited.")
	return nil
}

func (sv *Server) Shutdown() {
	log := util.GetLogger()
	sv.shutdownFlag = true
	log.Info().Msg("Shutdown initiated")
}

// Socket handling thread
func (sv *Server) awaitListener() {
	log := util.GetLogger()
	log.With().Str("module", "awaitListener")
	log.Info().Msg("Start Listener")

	for {
		log.Debug().Msg("Wait for client")

		// Accept() block execution.
		// continue when new client accepted or listener is closed (=shutdown)
		conn, err := sv.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Error().Err(err).Msg("Stop accept because of shutdown")
			} else {
				// Only network error. dont shutdown server
				log.Error().Err(err).Msg("Accept failed. continue")
				continue
			}
		}

		if sv.shutdownFlag {
			break
		}

		log.Debug().Msg("Accepted new client")
		go sv.dispatch(conn)

	}

	sv.procmanCh <- procman.RES_SHUTDOWN_DONE
	log.Info().Msg("Socket server stopped.")
}

// Connected socket handling thread
// move to module
func (sv *Server) dispatch(conn net.Conn) {
	log := util.GetLogger()

	log.Debug().Msg("request_dispatcher")
	scanner := bufio.NewScanner(conn)
	logchannel := make(chan string, 100)
	stopChannel := make(chan string, 1)

	start := time.Now()
	for scanner.Scan() {
		line := scanner.Text() // スキャンした内容を文字列で取得
		lineStr := strings.ToUpper(strings.TrimSpace(string(line)))

		log.Debug().Str("set-type", lineStr).Msg("Received")
		if lineStr == "LISTENER" {
			log.Debug().Msg("Start listener")
			go logcollecter.LogListener(conn, nil, stopChannel, logchannel)
			go logcollecter.PseudoLogSender(nil, stopChannel, logchannel)
		} else if lineStr == "COMMANDER" {
			log.Debug().Msg("Start commander")
			// go job.RequestHandler(conn, shutdownChannel, stopChannel, logchannel)
			// TODO
		}

		if time.Since(start).Seconds() > 10 {
			log.Error().Msg("Timeout waiting first message. Closing connection")
			conn.Close()
			break
		}
	}
}
