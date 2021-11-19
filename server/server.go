package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/pkg/server"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Server struct {
	procman.ProcmanModuleStruct
	listener net.Listener
}

func (sv *Server) IsInitialized() bool {
	return sv.Initialized
}

func (sv *Server) Initialize() error {
	sv.Name = "Server"
	return nil
}

func (sv *Server) GetName() string {
	return sv.Name
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

func (sv *Server) Start(inCh <-chan string, outCh chan<- string) error {
	sv.FromProcmanCh = inCh
	sv.ToProcmanCh = outCh
	log := util.GetLogger()

	log.Info().Msg("Starting socket server.")
	sv.ShutdownFlag = false

	err := sv.startListen()
	if err != nil {
		return err
	}
	log.Debug().Msg("TCP Socket start")

	go sv.awaitListener()
	log.Info().Msg("Socket server started.")

	for {
		select {
		case v := <-sv.FromProcmanCh:
			log.Debug().Msgf("Got request from procman %s", v)
		default:
		}

		if sv.ShutdownFlag {
			break
		}

		time.Sleep(procman.MAIN_LOOP_WAIT)
	}

	sv.listener.Close()
	log.Debug().Msg("Main thread exited.")
	return nil
}

func (sv *Server) Shutdown() {
	log := util.GetLogger()
	sv.ShutdownFlag = true
	log.Debug().Msg("Shutdown initiated")
}

// Socket handling thread
func (sv *Server) awaitListener() {
	log := util.GetLogger()
	log.With().Str("module", "awaitListener")
	log.Info().Msg("Start Listener")

	sv.ToProcmanCh <- procman.RES_STARTUP_DONE

	for {
		log.Debug().Msg("Wait for client")

		// Accept() block execution.
		// continue when new client accepted or listener is closed (=shutdown)
		conn, err := sv.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Debug().Err(err).Msg("Stop accept because of shutdown")
			} else {
				// Only network error. dont shutdown server
				log.Error().Err(err).Msg("Accept failed. continue")
				continue
			}
		}

		if sv.ShutdownFlag {
			break
		}

		log.Debug().Msg("Accepted new client")
		go sv.dispatch(conn)

	}

	sv.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	log.Info().Msg("Socket server stopped.")
}

// Connected socket handling thread
// move to module
func (sv *Server) dispatch(conn net.Conn) {
	log := util.GetLoggerWithSource("Dispatch")

	log.Debug().Msg("request_dispatcher")
	scanner := bufio.NewScanner(conn)
	logchannel := make(chan string, 100)
	stopChannel := make(chan string, 1)

	start := time.Now()
	for scanner.Scan() {
		line := scanner.Text() // スキャンした内容を文字列で取得
		lineStr := strings.ToUpper(strings.TrimSpace(string(line)))

		log.Debug().Str("set-type", lineStr).Msg("Received from client")
		switch lineStr {
		case "LISTENER":
			log.Debug().Msg("Start listener")
			go logcollecter.LogListener(conn, nil, stopChannel, logchannel)
			go logcollecter.PseudoLogSender(nil, stopChannel, logchannel)
		case "COMMANDER":
			log.Debug().Msg("Start commander")
		case "PING":
			res := server.ResPong{}
			res.Message = "PONG"
			msgBytes, _ := json.Marshal(res)
			msg := string(msgBytes) + "\n" + server.MSG_SEPARATOR

			sentBytes, err := io.Copy(conn, bytes.NewBufferString(msg))
			if err != nil {
				log.Err(err).Msg("Send PONG response failed")
			}
			fmt.Println(bytes.NewBufferString(msg).Len())
			log.Debug().Int64("bytes", sentBytes).Msg("Sent")
			conn.Close()
			return
		}

		if time.Since(start).Seconds() > 10 {
			log.Error().Msg("Timeout waiting first message. Closing connection")
			conn.Close()
			break
		}
	}
}
