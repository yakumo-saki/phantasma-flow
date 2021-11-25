package main

import (
	"fmt"
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
	"github.com/yakumo-saki/phantasma-flow/metrics"
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
	defDir := os.Getenv("PHFLOW_DEF_DIR")
	dataDir := os.Getenv("PHFLOW_DATA_DIR")
	tempDir := os.Getenv("PHFLOW_TEMP_DIR")
	if homeDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("Get homedir fail, Please set PHFLOW_HOME environment value.")
		}
		homeDir = path.Join(home, ".config", "phantasma-flow")
	}
	if defDir == "" {
		defDir = path.Join(homeDir, "definitions")
	}
	if dataDir == "" {
		dataDir = path.Join(homeDir, "data")
	}
	if tempDir == "" {
		tempDir = path.Join(homeDir, "temp")
	}

	isNotGoodDir(homeDir, "PHFLOW_HOME_DIR")
	isNotGoodDir(defDir, "PHFLOW_DEF_DIR")
	isNotGoodDir(dataDir, "PHFLOW_DATA_DIR")
	isNotGoodDir(tempDir, "PHFLOW_TEMP_DIR")

	makeSureHomeDirExists(defDir, dataDir, tempDir)

	return homeDir

}

func isNotGoodDir(dirname string, name string) {
	fmt.Println(dirname, name)
	if dirname != "" {
		st, err := os.Stat(dirname)
		if os.IsNotExist(err) {
			// not exist is ok. try to create after this
			return
		}
		if st == nil {
			return
		}
		if !st.IsDir() {
			panic(name + " is defined but it is file. It must be directory:" + dirname)
		}
	}

}

func makeSureHomeDirExists(defDir, dataDir, tempDir string) {
	mkdir := func(p string) {
		if e := os.MkdirAll(p, 0750); e != nil {
			panic(fmt.Sprintf("mkdir failed %s %e\n", p, e))
		}
	}

	dCfg := path.Join(defDir, "config")
	dJob := path.Join(defDir, "job")
	dNode := path.Join(defDir, "node")

	mkdir(dCfg)
	mkdir(dJob)
	mkdir(dNode)

	dlog := path.Join(dataDir, "log")
	dmeta := path.Join(dataDir, "meta")

	mkdir(dlog)
	mkdir(dmeta)

	mkdir(tempDir)
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
	processManager.AddService(&metrics.PrometeusExporterModule{})

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
