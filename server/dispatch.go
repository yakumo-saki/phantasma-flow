package server

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type server struct {
	globalCh     chan string
	shutdownFlag bool
	listener     net.Listener
}

var srv server

func (sv *server) startListen() error {
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

func (sv *server) start() {
	log := util.GetLogger()

	log.Info().Msg("Starting socket server.")
	sv.shutdownFlag = false

	go sv.awaitListener(sv.globalCh)
	log.Debug().Msg("TCP Socket start")

	log.Info().Msg("Socket server started.")

	for {
		select {
		case v := <-sv.globalCh:
			log.Info().Msgf("Got shutdown message %s", v)
			sv.listener.Close()
			log.Info().Msg("Socket closed.")
			sv.shutdownFlag = true
		default:
		}
		if sv.shutdownFlag {
			break
		}
	}
	log.Info().Msg("Socket server stopped.")
}

func (sv *server) awaitListener(globalCh <-chan string) {
	log := util.GetLogger()
	log.With().Str("module", "awaitListener")
	log.Info().Msg("Start Listener")
	for {

		log.Debug().Msg("Wait for client")
		conn, err := sv.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Error().Err(err).Msg("Stop accept because of shutdown")
				return
			} else {
				// Only network error. dont shutdown server
				log.Error().Err(err).Msg("Accept failed. continue")
				continue
			}
		}

		log.Debug().Msg("Accepted new client")
		go sv.dispatch(conn, globalCh)

	}
}

func (sv *server) dispatch(conn net.Conn, shutdownChannel <-chan string) {
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
			go logcollecter.LogListener(conn, shutdownChannel, stopChannel, logchannel)
			go logcollecter.PseudoLogSender(shutdownChannel, stopChannel, logchannel)
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

func Initialize(globalChannel chan string) error {
	srv.globalCh = globalChannel
	return nil
}

func Start() error {
	err := srv.startListen()
	if err != nil {
		return err
	}
	go srv.start()
	return nil
}
