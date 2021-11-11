package main

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/yakumo-saki/phantasma-flow/dispatcher"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func main() {
	log := util.GetLogger()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	psock, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Error().Err(err).Msg("Failed listen tcp.")
		return
	}

	allConn := &[]net.Conn{}
	shutdownChannel := make(chan string, 1)

	go awaitListener(shutdownChannel, psock)
	log.Debug().Msg("TCP Socket start")

	shutdownFlag := false
	for {
		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Got stop signal")
			psock.Close()
			log.Info().Msg("Socket closed.")
			log.Info().Msg("Awaiting shutdown of other threads.")
			shutdownChannel <- "SHUTDOWN"
			shutdownFlag = true
			log.Info().Msg("Await done. Shutdown.")
		default:
		}
		if shutdownFlag {
			break
		}
	}

	for _, c := range *allConn {
		c.Close()
	}
}

func awaitListener(shutdown <-chan string, psock net.Listener) {
	log := util.GetLogger()
	log.With().Str("module", "awaitListener")
	log.Info().Msg("Start Listener")
	for {

		log.Debug().Msg("Wait for client")
		conn, err := psock.Accept()
		log.Debug().Msg("Accepted client")
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
		go dispatcher.Dispatch(conn, shutdown)

	}
}
