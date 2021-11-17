package procman

type ProcmanModuleStruct struct {
	ProcmanCh    chan string
	ShutdownFlag bool
	Name         string // Recommended for debug
	Initialized  bool
}

type ProcmanModule interface {
	IsInitialized() bool
	Initialize() error
	Start(chan string) error
	Shutdown()
	GetName() string
}

const RES_SHUTDOWN_DONE = "SHUTDOWN_DONE" // shutdown done

const REQ_PAUSE = "PAUSE" // not used yet.
