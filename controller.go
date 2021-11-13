package main

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/yakumo-saki/phantasma-flow/controller"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// TODO: get path from something
// ENV or bootstrap parameter
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("get homedir fail")
	}

	return path.Join(home, "phantasma-flow")

}

func main() {
	log := util.GetLogger()

	log.Info().Msg("Starting Phantasma flow version 0.0.0")

	// at first Initialize repository for all configs
	cfgpath := getConfigPath()
	err := repository.Initialize(cfgpath)
	if err != nil {
		log.Err(err).Msg("Error occured in reading initialize data")
		return
	}

	// Start modules
	shutdownChannel := make(chan string, 1)
	startServer(shutdownChannel)
}

func startServer(shutdown chan string) {
	log := util.GetLogger()

	log.Info().Msg("Starting socket server.")

	// Finally start listening
	// TODO change port by config
	psock, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Error().Err(err).Msg("Failed listen tcp.")
		return
	}

	log.Info().Msg("Starting signal handling.")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go awaitListener(shutdown, psock)
	log.Debug().Msg("TCP Socket start")

	shutdownFlag := false
	for {
		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Got stop signal")
			psock.Close()
			log.Info().Msg("Socket closed.")
			log.Info().Msg("Awaiting shutdown of other threads.")
			shutdown <- "SHUTDOWN"
			shutdownFlag = true
			log.Info().Msg("Await done. Shutdown.")
		default:
		}
		if shutdownFlag {
			break
		}
	}
}

func awaitListener(shutdown <-chan string, psock net.Listener) {
	log := util.GetLogger()
	log.With().Str("module", "awaitListener")
	log.Info().Msg("Start Listener")
	for {

		log.Debug().Msg("Wait for client")
		conn, err := psock.Accept()
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
		go controller.Dispatch(conn, shutdown)

	}
}
