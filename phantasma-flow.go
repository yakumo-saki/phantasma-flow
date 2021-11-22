package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/yakumo-saki/phantasma-flow/global"
	"github.com/yakumo-saki/phantasma-flow/jobscheduler"
	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
	"github.com/yakumo-saki/phantasma-flow/node"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/procmanExample"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/server"
	"github.com/yakumo-saki/phantasma-flow/util"
)

const DEBUG = false
const myname = "main"

// Get phantasma-flow home path.
// ENV or ~/.config/phantasma-flow
func getHomeDir() string {
	util.GetLoggerWithSource(myname, "homedir")
	homeDir := os.Getenv("PHFLOW_HOME")
	if homeDir != "" {
		st, err := os.Stat(homeDir)
		if !st.IsDir() {

			panic("PHFLOW_HOME is defined but it is file. It must be directory:" + homeDir)
		}
		if os.IsNotExist(err) {
			// not exist is ok. try to create after this
		} else {
			panic("can't stat PHFLOW_HOME:" + homeDir)
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic("get homedir fail")
	}
	homeDir = path.Join(home, "phantasma-flow")
	makeSureHomeDirExists(homeDir)

	return homeDir

}

func main() {
	log := util.GetLoggerWithSource("main")

	log.Info().Msgf("Starting Phantasma flow version %s (commit %s) %s",
		global.VERSION, global.COMMIT, global.URL)

	// at first Initialize repository for all configs
	repo := repository.Repository{}
	cfgpath := getHomeDir()
	err := repo.Initialize(cfgpath)
	if err != nil {
		log.Error().Err(err).Msg("Error occured at reading initialize data")
		log.Error().Msg("Maybe data is corrupted or misseditted.")
		return
	}

	// Start modules
	hub := messagehub_impl.MessageHub{}
	messagehub.SetMessageHub(&hub)
	hub.Initialize()
	hub.StartSender()

	procmanCh := make(chan string, 1) // controller to processManager. signal only
	processManager := procman.NewProcessManager(procmanCh)

	processManager.Add(&procmanExample.MinimalProcmanModule{})
	processManager.AddService(&logcollecter.LogListenerModule{})
	processManager.AddService(&jobscheduler.JobScheduler{})
	processManager.AddService(&node.NodeManager{})

	processManager.Start()

	log.Info().Msg("Starting signal handling.")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Load definitions
	repo.SendAllNodes() // must send node before job (must exist node, job requires)
	repo.SendAllJobs()

	waitForMessageHub(&log, &hub)
	// XXX: ノードとかジョブが行き渡ったことを確認する必要がある？
	// nodeDef とか JobDef を送った数の分のノードができたことをチェックする？

	// main loop
	processManager.AddService(&server.Server{})
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
	shutdownFlag := false
	for {
		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Got signal")
			shutdownFlag = true
			shutdown(&processManager, &hub)
		case <-debugCh:
			log.Warn().Msg("Debug shutdown start.")
			shutdownFlag = true
			shutdown(&processManager, &hub)
		default:
		}

		if shutdownFlag {
			break
		}

		time.Sleep(procman.MAIN_LOOP_WAIT)
	}
	log.Info().Msg("Phantasma-flow stopped.")
}

func shutdown(pm *procman.ProcessManager, hub *messagehub_impl.MessageHub) {
	log := util.GetLoggerWithSource("shutdown")
	hub.Shutdown()
	r1, r2 := pm.Shutdown()
	log.Info().Str("modules", r1).Str("services", r2).Msg("Threads shutdown done.")
}

func waitForMessageHub(log *zerolog.Logger, hub *messagehub_impl.MessageHub) {
	for {
		if hub.GetQueueLength() == 0 {
			log.Debug().Msg("Wait for message hub done.")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
