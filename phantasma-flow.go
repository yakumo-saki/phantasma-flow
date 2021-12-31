package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/global"
	"github.com/yakumo-saki/phantasma-flow/job/executer"
	"github.com/yakumo-saki/phantasma-flow/job/nodemanager"
	"github.com/yakumo-saki/phantasma-flow/job/scheduler"
	"github.com/yakumo-saki/phantasma-flow/logexporter/logfileexporter"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
	"github.com/yakumo-saki/phantasma-flow/metalog"
	"github.com/yakumo-saki/phantasma-flow/metrics"
	"github.com/yakumo-saki/phantasma-flow/pprofserver"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/procmanExample"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/repository/validater"
	"github.com/yakumo-saki/phantasma-flow/server"
	"github.com/yakumo-saki/phantasma-flow/util"
)

const DEBUG = false

func main() {
	log := util.GetLoggerWithSource("main")

	log.Info().Msgf("Starting Phantasma flow version %s (commit %s) %s",
		global.VERSION, global.COMMIT, global.URL)

	// at first Initialize repository for all configs
	repo := startRepository()

	// Start modules
	hub := startMessageHub()

	execute := executer.GetInstance()
	logMetaMan := metalog.GetInstance()

	procmanCh := make(chan string, 1) // controller to processManager. signal only
	processManager := procman.NewProcessManager(procmanCh)

	processManager.Add(&procmanExample.MinimalProcmanModule{})
	processManager.Add(&metrics.PrometeusExporterModule{})
	processManager.Add(&server.Server{})
	processManager.Add(&pprofserver.PprofServer{})
	processManager.AddService(10, &logfileexporter.LogFileExporter{})
	processManager.AddService(11, logMetaMan)
	processManager.AddService(80, nodemanager.GetInstance())
	processManager.AddService(90, execute)
	processManager.AddService(91, &scheduler.JobScheduler{})

	processManager.Start()

	validater.ValidateAllJobDef()

	// Load definitions
	repo.SendAllNodes() // must send node before job (must exist node, job requires)
	messagehub.WaitForQueueEmpty("Wait for node registration")
	repo.SendAllJobs()
	messagehub.WaitForQueueEmpty("Wait for job registration")

	log.Info().Msg("Starting signal handling.")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Phantasma-flow started.")

	// for debug
	debugCh := make(chan string, 1)
	if DEBUG {
		go func() {
			log := util.GetLogger()
			for i := 0; i < 3; i++ {
				time.Sleep(1 * time.Second)
				log.Debug().Msgf("Wait for timeout %d", i)
			}
			debugCh <- "SHUTDOWN"
		}()
	}

	// wait for stop signal
	for {
		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Got signal")
			shutdownProcMan(&processManager, hub)
			goto shutdown
		case <-debugCh:
			log.Warn().Msg("Debug shutdown start.")
			shutdownProcMan(&processManager, hub)
			goto shutdown
		}
	}

shutdown:
	log.Info().Msg("Phantasma-flow stopped.")
}

func startRepository() *repository.Repository {
	repo := repository.GetRepository()
	err := repo.Initialize()
	if err != nil {
		log.Error().Err(err).Msg("Error occured at reading initialize data")
		log.Error().Msg("Maybe data is corrupted or misseditted.")
		return nil
	}

	return repo
}

// StartMessageHub
func startMessageHub() *messagehub_impl.MessageHub {
	hub := messagehub_impl.MessageHub{}
	messagehub.SetMessageHub(&hub)
	hub.Initialize()
	hub.StartSender()
	return &hub
}

func shutdownProcMan(pm *procman.ProcessManager, hub *messagehub_impl.MessageHub) {
	log := util.GetLoggerWithSource("shutdown")
	hub.Shutdown()
	r1, r2 := pm.Shutdown()
	log.Info().Str("modules", r1).Str("services", r2).Msg("Threads shutdown done.")
}
