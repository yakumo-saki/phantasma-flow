package procman

type ProcmanModule interface {
	Initialize(chan string) error
	Start() error
	Shutdown()
	GetName() string
}

const RES_SHUTDOWN_DONE = "SHUTDOWN_DONE" // shutdown done

const REQ_PAUSE = "PAUSE" // not used yet.
