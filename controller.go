package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/procmanExample"
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
	procmanCh := make(chan string, 1) // controller to processManager. signal only
	processManager := procman.NewProcessManager(procmanCh)

	processManager.Add(&procmanExample.MinimalProcmanModule{})
	processManager.AddService(&server.Server{})
	processManager.AddService(&logcollecter.LogListenerModule{})

	go processManager.Start()

	log.Info().Msg("Starting signal handling.")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// for debug
	debugCh := make(chan string, 1)
	go func() {
		log := util.GetLogger()
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			log.Debug().Msgf("Wait for timeout %d", i)
		}
		debugCh <- "SHUTDOWN"
	}()

	shutdownFlag := false
	for {
		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Got signal")
			r1, r2 := processManager.Shutdown()
			shutdownFlag = true
			log.Info().Str("modules", r1).Str("services", r2).Msg("Threads shutdown done.")
		case <-debugCh:
			log.Warn().Msg("Debug shutdown start.")
			r1, r2 := processManager.Shutdown()
			shutdownFlag = true
			log.Info().Str("modules", r1).Str("services", r2).Msg("Threads shutdown done.")
		default:
		}

		if shutdownFlag {
			break
		}

		time.Sleep(procman.MAIN_LOOP_WAIT)
	}
}
