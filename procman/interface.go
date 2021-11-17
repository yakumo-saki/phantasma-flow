package procman

type ProcmanModuleStruct struct {
	FromProcmanCh <-chan string
	ToProcmanCh   chan<- string
	ShutdownFlag  bool
	Name          string // Recommended for debug
	Initialized   bool
}

type ProcmanModule interface {
	IsInitialized() bool
	Initialize() error
	Start(<-chan string, chan<- string) error
	Shutdown()
	GetName() string
}

// Module -> Procman
const RES_STARTUP_DONE = "STARTUP_DONE"   // response: Start() done
const RES_SHUTDOWN_DONE = "SHUTDOWN_DONE" // response: Shutdown() done

// Procman -> module
const REQ_PAUSE = "PAUSE" // not used yet.
