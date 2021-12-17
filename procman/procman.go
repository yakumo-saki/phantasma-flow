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
	success := p.addImpl(TYPE_MOD, p.workerModules, module, 0)
	if !success {
		panic("Add failed. name=" + module.GetName())
	}
}

// Add module as Service module
// if procman is started, module is automaticaly start
func (p *ProcessManager) AddService(order uint8, module ProcmanModule) {
	success := p.addImpl(TYPE_SVC, p.serviceModules, module, order)
	if !success {
		panic("AddService failed. name=" + module.GetName())
	}
}

//
func (p *ProcessManager) addImpl(typeName string, modmap map[string]*process, module ProcmanModule, order uint8) bool {
	log := util.GetLoggerWithSource(myname, "add")

	toCh := make(chan string, 1)
	fromCh := make(chan string, 1)

	if !module.IsInitialized() {
		module.Initialize()
	}

	name := module.GetName()
	if typeName == TYPE_SVC {
		name = fmt.Sprintf("%03v %s", order, module.GetName())
	}
	if name == "" {
		msg := fmt.Sprintf("[%s] empty name is not allowed", typeName)
		panic(msg)
	}

	_, ok := modmap[name]
	if ok {
		msg := fmt.Sprintf("[%s] name %s is already registered.", typeName, name)
		panic(msg)
	}

	proc := process{toModCh: toCh, fromModCh: fromCh, module: module}
	modmap[name] = &proc

	// Automatic start modules, when procman is already started
	if p.startupDone {
		result := p.StartImpl(typeName, modmap)
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
	svcResult := p.StartImpl(TYPE_SVC, p.serviceModules)
	if svcResult == REASON_COMPLETE {
		log.Debug().Msgf("[%s] All services started", TYPE_SVC)
	} else {
		msg := fmt.Sprintf("[%s] Some or all services failed to start. %s", TYPE_SVC, svcResult)
		log.Error().Msgf(msg)
		panic(msg)
	}

	modResult := p.StartImpl(TYPE_MOD, p.workerModules)
	if modResult == REASON_COMPLETE {
		log.Debug().Msgf("[%s] All modules started", TYPE_MOD)
	} else {
		msg := fmt.Sprintf("[%s] Some or all modules failed to start. %s", TYPE_SVC, modResult)
		log.Error().Msgf(msg)
		panic(msg)
	}

	p.startupDone = true
}

func (p *ProcessManager) StartImpl(typeName string, modmap map[string]*process) string {
	var result string
	switch typeName {
	case TYPE_MOD:
		result = p.startModules(modmap)
	case TYPE_SVC:
		result = p.startServices(modmap)
	default:
		panic(fmt.Sprintf("Unknown typeName %s", typeName))
	}
	return result
}

func (p *ProcessManager) Shutdown() (string, string) {
	// log := util.GetLogger()

	reason1 := p.stopImpl(TYPE_MOD, p.workerModules)

	// todo allow timeout?

	reason2 := p.stopImpl(TYPE_SVC, p.serviceModules)

	p.shutdownInitiated = true

	return reason1, reason2
}

func (p *ProcessManager) stopImpl(typeName string, modmap map[string]*process) string {
	var result string
	switch typeName {
	case TYPE_MOD:
		result = p.stopModules(modmap)
	case TYPE_SVC:
		result = p.stopServices(modmap)
	default:
		panic(fmt.Sprintf("Unknown typeName %s", typeName))
	}
	return result
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
