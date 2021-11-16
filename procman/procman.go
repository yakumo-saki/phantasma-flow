package procman

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

type process struct {
	module   ProcmanModule
	channel  chan string
	started  bool
	shutdown bool // shutdown complete flag of Procmanager
}

type ProcessManager struct {
	modules        map[string]process // start last, shutdown first
	serviceModules map[string]process // start first, shutdown last
	inChannel      chan string
	shutdownFlag   bool // shutdown initiated flag of Procmanager
}

const MSG_SHUTDOWN_COMPLETE = "SHUTDOWN COMPLETE"

func (p *ProcessManager) Add(module ProcmanModule) bool {
	return p.AddImpl("Module", p.modules, module)
}

func (p *ProcessManager) AddService(module ProcmanModule) bool {
	return p.AddImpl("Service", p.serviceModules, module)
}

func (p *ProcessManager) AddImpl(typeName string, modmap map[string]process, module ProcmanModule) bool {
	log := util.GetLogger()

	name := module.GetName()
	if _, ok := p.modules[name]; ok {
		log.Error().Msgf("[%s] name %s is already subscribed.", typeName, name)
		return false
	}

	proc := process{channel: make(chan string, 1), module: module}

	modmap[name] = proc
	return true
}

func (p *ProcessManager) Start() {
	for _, proc := range p.serviceModules {
		proc.module.Initialize(proc.channel)
		proc.module.Start()
		proc.started = true
	}

	for _, proc := range p.modules {
		proc.module.Initialize(proc.channel)
		proc.module.Start()
		proc.started = true
	}

}

func (p *ProcessManager) Shutdown() (string, string) {
	// log := util.GetLogger()

	reason1 := p.shutdownImpl("modules", p.modules)

	//todo

	reason2 := p.shutdownImpl("services", p.serviceModules)

	p.shutdownFlag = true

	return reason1, reason2
}

func (p *ProcessManager) shutdownImpl(typeName string, modmap map[string]process) string {
	log := util.GetLogger()

	var reason string

	if len(modmap) == 0 {
		log.Debug().Msgf("[%s] Has no modules.", typeName)
		return RES_SHUTDOWN_DONE
	}

	timeoutCh := make(chan string, 1)
	go func() {
		time.Sleep(10 * time.Second)
		log.Debug().Msgf("[%s] Shutdown timeout reached", typeName)
		timeoutCh <- "TIMEOUT"
	}()

	for k, proc := range modmap {
		log.Debug().Msgf("[%s] Sending shutdown request to %s", typeName, k)
		proc.shutdown = false
		proc.module.Shutdown()
	}

	for {
		stop := false
		for k, proc := range modmap {
			select {
			case v := <-proc.channel:
				if v == RES_SHUTDOWN_DONE {
					proc.shutdown = true
					if p.isShutdownComplete(modmap) {
						stop = true
						reason = MSG_SHUTDOWN_COMPLETE
					}
				} else {
					log.Warn().Str("module", k).Str("message", v).Msg("Unexpected response")
				}
			case <-timeoutCh:
				reason = "TIMEOUT"
				p.outputTimeoutLog(typeName, p.modules)

				stop = true
			default:
			}
		}

		if stop {
			log.Debug().Msgf("Stop for %s", reason)
			break
		}
	}
	return reason
}

func (p *ProcessManager) isShutdownComplete(modmap map[string]process) bool {
	for _, proc := range modmap {
		if !proc.shutdown {
			return false
		}
	}
	return true
}

func (p *ProcessManager) outputTimeoutLog(typeName string, modmap map[string]process) {
	log := util.GetLogger()

	for name, proc := range modmap {
		if !proc.shutdown {
			log.Error().Str("name", name).
				Msgf("[%s]Do not shutdown complete until timeout.", typeName)
		}
	}

}

func NewProcessManager(channel chan string) ProcessManager {
	var p ProcessManager
	p.inChannel = channel
	p.shutdownFlag = false
	p.modules = make(map[string]process, 10)
	p.serviceModules = make(map[string]process, 10)

	return p
}
