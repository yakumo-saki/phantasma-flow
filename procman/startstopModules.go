package procman

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// StartModules
func (p *ProcessManager) startModules(modmap map[string]*process) string {
	log := util.GetLoggerWithSource(myname, "start")

	if len(modmap) == 0 {
		return REASON_COMPLETE
	}

	typeName := TYPE_MOD

	// start modules/services
	for _, proc := range modmap {
		if !proc.started {
			go proc.module.Start(proc.toModCh, proc.fromModCh)
			log.Debug().Msgf("[%s] Request starting %s.", typeName, proc.module.GetName())
		}
	}

	reason := "UNKNOWN"
	for {
		for name, proc := range modmap {
			select {
			case v := <-proc.fromModCh:
				if v == RES_STARTUP_DONE {
					log.Debug().Msgf("[%s] %s is started", typeName, name)
					proc.started = true
					proc.shutdown = false
					if p.isStartupComplete(modmap) {
						return REASON_COMPLETE
					}
				} else {
					log.Warn().Str("module", name).Str("message", v).Msg("Unexpected response")
				}
			case <-time.After(15 * time.Second):
				reason = REASON_TIMEOUT
				p.outputTimeoutLog(typeName, "startup", p.workerModules)
				log.Fatal().Msgf("[%s] Startup timeout reached. '%s' is not started.", typeName, name)
				return reason
			default:
			}
		}
	}
}

// Stop modules.
func (p *ProcessManager) stopModules(modmap map[string]*process) string {
	log := util.GetLoggerWithSource(myname, "shutdown")

	var reason string

	typeName := TYPE_MOD

	if len(modmap) == 0 {
		log.Debug().Msgf("[%s] Has no modules.", typeName)
		return REASON_COMPLETE
	}

	for k, proc := range modmap {
		log.Debug().Msgf("[%s] Sending shutdown request to %s", typeName, k)
		proc.shutdown = false
		proc.module.Shutdown()
	}

	for {
		stop := false
		for k, proc := range modmap {
			select {
			case v := <-proc.fromModCh:
				if v == RES_SHUTDOWN_DONE {
					proc.shutdown = true
					if p.isShutdownComplete(modmap) {
						stop = true
						reason = REASON_COMPLETE
						break
					}
				} else {
					log.Warn().Str("module", k).Str("message", v).Msg("Unexpected response")
				}
			case <-time.After(15 * time.Second):
				reason = "TIMEOUT"
				stop = true
			default:
			}
		}

		if stop {
			if reason == "TIMEOUT" {
				p.outputTimeoutLog(typeName, "shutdown", modmap)
			}

			log.Debug().Str("cause", reason).Msgf("[%s] Shutdown done.", typeName)
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
	return reason
}
