package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/server"
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
		log.Err(err).Msg("Error occured at reading initialize data")
		return
	}

	// Start modules

	globalCh := make(chan string, 1)
	procman := procman.NewProcessManager(globalCh)
	ch, ok := procman.Subscribe("server")
	if !ok {
		panic("Duplicate name of processManager")
	}

	server.Initialize(ch)
	err = server.Start()
	if err != nil {
		log.Err(err).Msg("Error occured at starting server")
		return
	}

	log.Info().Msg("Starting signal handling.")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	shutdownFlag := false
	for {
		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Got signal")
			r := procman.Shutdown()
			shutdownFlag = true
			log.Info().Str("result", r).Msg("Threads shutdown done.")
		default:
		}

		if shutdownFlag {
			break
		}
	}
}
