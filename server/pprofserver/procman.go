package pprofserver

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type PprofServer struct {
	procman.ProcmanModuleStruct
}

func (sv *PprofServer) IsInitialized() bool {
	return sv.Initialized
}

func (sv *PprofServer) Initialize() error {
	sv.Initialized = true
	return nil
}

func (sv *PprofServer) GetName() string {
	return "PprofServer"
}
func (sv *PprofServer) Start(inCh <-chan string, outCh chan<- string) error {
	const NAME = "main"
	log := util.GetLoggerWithSource(sv.GetName(), NAME)

	sv.FromProcmanCh = inCh
	sv.ToProcmanCh = outCh

	sv.ToProcmanCh <- procman.RES_STARTUP_DONE

	// Get setting from repository
	cfg := sv.GetConfig()
	if !cfg.Enabled {
		log.Debug().Msgf("Pprof server for debugging disabled. (Recommended for production)")
		return nil
	}

	go func() {
		log.Info().Msgf("Debug interface listen on %s", cfg.ListenAddrAndPort)
		log.Info().Msgf("%v", http.ListenAndServe(cfg.ListenAddrAndPort, nil))
	}()

	return nil
}

func (sv *PprofServer) GetConfig() objects.PprofServerConfig {
	bareCfg := repository.GetRepository().GetConfigByKind(objects.KIND_PPROF_SERVER_CFG)
	cfg := objects.PprofServerConfig{}
	if bareCfg != nil {
		cfg = bareCfg.(objects.PprofServerConfig)
		return cfg
	}

	cfg.Kind = objects.KIND_PPROF_SERVER_CFG
	cfg.ListenAddrAndPort = "localhost:6060"
	cfg.Enabled = false
	return cfg
}

func (sv *PprofServer) Shutdown() {
	sv.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
}
