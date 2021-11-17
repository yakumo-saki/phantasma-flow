package procman

import (
	"fmt"
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

type process struct {
	module    ProcmanModule
	toModCh   chan string
	fromModCh chan string
	started   bool // startup done complete flag
	shutdown  bool // shutdown done complete flag
}

type ProcessManager struct {
	workerModules  map[string]*process // start last, shutdown first
	serviceModules map[string]*process // start first, shutdown last
	inChannel      chan string
	shutdownFlag   bool // shutdown initiated flag of Procmanager
	startupFlag    bool // Startup done flag
}

const MAIN_LOOP_WAIT = 1000 * time.Millisecond // Recommended wait for message loop

const REASON_COMPLETE = "COMPLETE"
const REASON_TIMEOUT = "TIMEOUT"

const TYPE_MOD = "modules"
const TYPE_SVC = "services"

func (p *ProcessManager) Add(module ProcmanModule) {
	success := p.AddImpl("Module", p.workerModules, module)
	if !success {
		panic("Add failed. name=" + module.GetName())
	}
}

func (p *ProcessManager) AddService(module ProcmanModule) {
	success := p.AddImpl("Service", p.serviceModules, module)
	if !success {
		panic("AddService failed. name=" + module.GetName())
	}
}

func (p *ProcessManager) AddImpl(typeName string, modmap map[string]*process, module ProcmanModule) bool {
	log := util.GetLogger()

	toCh := make(chan string, 1)
	fromCh := make(chan string, 1)

	if !module.IsInitialized() {
		module.Initialize()
	}

	name := module.GetName()
	if name == "" {
		msg := fmt.Sprintf("[%s] empty name is not allowed", typeName)
		panic(msg)
	}

	_, ok := modmap[name]
	if ok {
		log.Error().Msgf("[%s] name %s is already registered.", typeName, name)
		return false
	}

	proc := process{toModCh: toCh, fromModCh: fromCh, module: module}
	modmap[name] = &proc

	log.Debug().Msgf("[%s] %s is added.", typeName, name)
	return true
}

// Blocks until all modules are start or not.
func (p *ProcessManager) Start() {
	log := util.GetLogger()
	p.startImpl(TYPE_SVC, p.serviceModules)
	log.Debug().Msgf("[%s] All services started", TYPE_SVC)

	p.startImpl(TYPE_MOD, p.workerModules)
	log.Debug().Msgf("[%s] All modules started", TYPE_MOD)
}

func (p *ProcessManager) startImpl(typeName string, modmap map[string]*process) string {
	log := util.GetLogger()
	for _, proc := range modmap {
		proc.module.Initialize()
		go proc.module.Start(proc.toModCh, proc.fromModCh)
		log.Debug().Msgf("[%s] Request starting %s .", typeName, proc.module.GetName())
	}

	timeoutCh := make(chan string, 1)
	go util.Timeout(timeoutCh, 10)

	reason := "UNKNOWN"
	for {
		for name, proc := range modmap {
			select {
			case v := <-proc.fromModCh:
				if v == RES_STARTUP_DONE {
					log.Debug().Msgf("[%s] %s is started", typeName, name)
					proc.started = true
					if p.isStartupComplete(modmap) {
						return REASON_COMPLETE
					}
				} else {
					log.Warn().Str("module", name).Str("message", v).Msg("Unexpected response")
				}
			case <-timeoutCh:
				reason = REASON_TIMEOUT
				p.outputTimeoutLog(typeName, "startup", p.workerModules)
				log.Debug().Msgf("[%s] Startup timeout reached", typeName)
				return reason
			default:
			}
		}
	}
}

func (p *ProcessManager) Shutdown() (string, string) {
	// log := util.GetLogger()

	reason1 := p.shutdownImpl("modules", p.workerModules)

	// todo allow timeout?

	reason2 := p.shutdownImpl("services", p.serviceModules)

	p.shutdownFlag = true

	return reason1, reason2
}

func (p *ProcessManager) shutdownImpl(typeName string, modmap map[string]*process) string {
	log := util.GetLogger()

	var reason string

	if len(modmap) == 0 {
		log.Debug().Msgf("[%s] Has no modules.", typeName)
		return REASON_COMPLETE
	}

	timeoutCh := make(chan string, 1)
	go func() {
		time.Sleep(10 * time.Second)
		timeoutCh <- "TIMEOUT"
		log.Debug().Msgf("[%s] Shutdown timeout reached", typeName)
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
			case v := <-proc.fromModCh:
				if v == RES_SHUTDOWN_DONE {
					proc.shutdown = true
					if p.isShutdownComplete(modmap) {
						stop = true
						reason = REASON_COMPLETE
					}
				} else {
					log.Warn().Str("module", k).Str("message", v).Msg("Unexpected response")
				}
			case <-timeoutCh:
				reason = "TIMEOUT"
				p.outputTimeoutLog(typeName, "shutdown", p.workerModules)

				stop = true
			default:
			}
		}

		if stop {
			log.Debug().Msgf("[%s] Stopped cause: %s", typeName, reason)
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
	return reason
}

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
	log := util.GetLogger()

	for name, proc := range modmap {
		if !proc.shutdown {
			log.Error().Str("name", name).
				Msgf("[%s]Do not response %s complete until timeout.", typeName, action)
		}
	}

}

func NewProcessManager(channel chan string) ProcessManager {
	var p ProcessManager
	p.inChannel = channel
	p.shutdownFlag = false
	p.startupFlag = false
	p.workerModules = make(map[string]*process, 10)
	p.serviceModules = make(map[string]*process, 10)

	return p
}
