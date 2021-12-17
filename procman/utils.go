package procman

import "github.com/yakumo-saki/phantasma-flow/util"

func (p *ProcessManager) isShutdownComplete(modmap map[string]*process) bool {
	for _, proc := range modmap {
		if !proc.shutdown {
			return false
		}
	}
	return true
}

func (p *ProcessManager) isStartupComplete(modmap map[string]*process) bool {
	for _, proc := range modmap {
		if !proc.started {
			return false
		}
	}
	return true
}

func (p *ProcessManager) outputTimeoutLog(typeName string, action string, modmap map[string]*process) {
	log := util.GetLoggerWithSource(myname, "timeout")

	for name, proc := range modmap {
		switch action {
		case "shutdown":
			if !proc.shutdown {
				log.Error().Str("name", name).
					Msgf("[%s] Do not response %s complete until timeout.", typeName, action)
			}
		default:
			if !proc.started {
				log.Error().Str("name", name).
					Msgf("[%s] Do not response %s complete until timeout.", typeName, action)
			}
		}
	}
}
