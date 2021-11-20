package procman

import (
	"context"
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
	workerModules     map[string]*process // start last, shutdown first
	serviceModules    map[string]*process // start first, shutdown last
	inChannel         chan string
	shutdownInitiated bool // shutdown initiated flag of Procmanager
	startupDone       bool // Startup done flag
}

const MAIN_LOOP_WAIT = 500 * time.Millisecond // Recommended wait for message loop

const REASON_COMPLETE = "COMPLETE"
const REASON_TIMEOUT = "TIMEOUT"

const TYPE_MOD = "modules"
const TYPE_SVC = "services"

const myname = "procman"

// Add module as Worker module
// if procman is started, module is automaticaly start
func (p *ProcessManager) Add(module ProcmanModule) {
	success := p.AddImpl("Module", p.workerModules, module)
	if !success {
		panic("Add failed. name=" + module.GetName())
	}
}

// Add module as Service module
// if procman is started, module is automaticaly start
func (p *ProcessManager) AddService(module ProcmanModule) {
	success := p.AddImpl("Service", p.serviceModules, module)
	if !success {
		panic("AddService failed. name=" + module.GetName())
	}
}

func (p *ProcessManager) AddImpl(typeName string, modmap map[string]*process, module ProcmanModule) bool {
	log := util.GetLoggerWithSource(myname, "add")

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

	// Automatic start modules, when procman is already started
	if p.startupDone {
		result := p.startImpl(typeName, modmap)
		if result == REASON_COMPLETE {
			log.Debug().Msgf("[%s] %s is added and started.", typeName, name)
		} else {
			log.Debug().Msgf("[%s] %s is added. but failed to start. %s", typeName, name, result)
			return false
		}
	} else {
		log.Debug().Msgf("[%s] %s is added.", typeName, name)
	}
	return true
}

// Blocks until all modules are start or not.
func (p *ProcessManager) Start() {
	log := util.GetLoggerWithSource(myname, "start")
	svcResult := p.startImpl(TYPE_SVC, p.serviceModules)
	if svcResult == REASON_COMPLETE {
		log.Debug().Msgf("[%s] All services started", TYPE_SVC)
	} else {
		msg := fmt.Sprintf("[%s] Some or all services failed to start. %s", TYPE_SVC, svcResult)
		log.Error().Msgf(msg)
		panic(msg)
	}

	modResult := p.startImpl(TYPE_MOD, p.workerModules)
	if modResult == REASON_COMPLETE {
		log.Debug().Msgf("[%s] All modules started", TYPE_MOD)
	} else {
		msg := fmt.Sprintf("[%s] Some or all modules failed to start. %s", TYPE_SVC, modResult)
		log.Error().Msgf(msg)
		panic(msg)
	}

	p.startupDone = true
}

func (p *ProcessManager) startImpl(typeName string, modmap map[string]*process) string {
	log := util.GetLoggerWithSource(myname, "start")
	for _, proc := range modmap {
		if !proc.started {
			go proc.module.Start(proc.toModCh, proc.fromModCh)
			log.Debug().Msgf("[%s] Request starting %s.", typeName, proc.module.GetName())
		}
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
					proc.shutdown = false
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

	p.shutdownInitiated = true

	return reason1, reason2
}

func (p *ProcessManager) shutdownImpl(typeName string, modmap map[string]*process) string {
	log := util.GetLoggerWithSource(myname, "shutdown")

	var reason string

	if len(modmap) == 0 {
		log.Debug().Msgf("[%s] Has no modules.", typeName)
		return REASON_COMPLETE
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
			case <-ctx.Done():
				reason = "TIMEOUT"
				stop = true
			default:
			}
		}

		if stop {
			if reason == "TIMEOUT" {
				p.outputTimeoutLog(typeName, "shutdown", modmap)
			}

			log.Debug().Msgf("[%s] Shutdown done. cause: %s", typeName, reason)
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

func NewProcessManager(channel chan string) ProcessManager {
	var p ProcessManager
	p.inChannel = channel
	p.shutdownInitiated = false
	p.startupDone = false
	p.workerModules = make(map[string]*process, 10)
	p.serviceModules = make(map[string]*process, 10)

	return p
}
