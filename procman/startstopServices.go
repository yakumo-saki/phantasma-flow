package procman

import (
	"sort"
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// startServices
func (p *ProcessManager) startServices(modmap map[string]*process) string {
	log := util.GetLoggerWithSource(myname, "start")

	typeName := TYPE_SVC

	if len(modmap) == 0 {
		return REASON_COMPLETE
	}

	// create order
	serviceNamesNotSorted := []string{}
	for name := range modmap {
		serviceNamesNotSorted = append(serviceNamesNotSorted, name)
	}
	serviceNames := sort.StringSlice(serviceNamesNotSorted)

	// each service have timeout limit.
	reason := "UNKNOWN"
	for _, name := range serviceNames {
		proc := modmap[name]

		if !proc.started {
			go proc.module.Start(proc.toModCh, proc.fromModCh)
			log.Debug().Msgf("[%s] Request starting %s.", typeName, proc.module.GetName())
		}

		select {
		case v := <-proc.fromModCh:
			if v == RES_STARTUP_DONE {
				log.Debug().Msgf("[%s] %s is started", typeName, name)
				proc.started = true
				proc.shutdown = false
			} else {
				log.Warn().Str("module", name).Str("message", v).Msg("Unexpected response")
			}
		case <-time.After(15 * time.Second):
			reason = REASON_TIMEOUT
			p.outputTimeoutLog(typeName, "startup", p.workerModules)
			log.Error().Msgf("[%s] Startup timeout reached. '%s' is not started.", typeName, name)
			return reason
		}
	}

	return REASON_COMPLETE
}

// Stop modules.
func (p *ProcessManager) stopServices(modmap map[string]*process) string {
	log := util.GetLoggerWithSource(myname, "shutdown")

	var reason string
	typeName := TYPE_SVC

	if len(modmap) == 0 {
		log.Debug().Msgf("[%s] Has no modules.", typeName)
		return REASON_COMPLETE
	}

	serviceNames := []string{}
	for name := range modmap {
		serviceNames = append(serviceNames, name)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(serviceNames)))

	for _, name := range serviceNames {
		proc := modmap[name]

		log.Debug().Msgf("[%s] Sending shutdown request to %s", typeName, name)
		proc.shutdown = false
		proc.module.Shutdown()

		select {
		case v := <-proc.fromModCh:
			if v == RES_SHUTDOWN_DONE {
				proc.shutdown = true
				if p.isShutdownComplete(modmap) {
					reason = REASON_COMPLETE
					break
				}
			} else {
				log.Warn().Str("module", name).Str("message", v).Msg("Unexpected response")
			}
		case <-time.After(15 * time.Second):
			reason = "TIMEOUT"
		}
	}

	if reason == "TIMEOUT" {
		p.outputTimeoutLog(typeName, "shutdown", modmap)
		return reason
	}

	log.Debug().Str("cause", reason).Msgf("[%s] Shutdown done.", typeName)

	return REASON_COMPLETE
}
