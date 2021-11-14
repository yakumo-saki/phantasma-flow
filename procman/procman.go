package procman

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

const SHUTDOWN_DONE = "SHUTDOWN_DONE" // shutdown done
const SHUTDOWN = "SHUTDOWN"           // shutdown request

type ProcessManager struct {
	channels     map[string]chan string
	inChannel    chan string
	shutdownFlag bool
}

func (p *ProcessManager) Subscribe(name string) (chan string, bool) {
	log := util.GetLogger()
	ch := make(chan string)

	if _, ok := p.channels[name]; ok {
		log.Error().Msgf("Name %s is already subscribed.", name)
		return nil, false
	}
	p.channels[name] = ch
	return ch, true
}

func (p *ProcessManager) Shutdown() string {
	log := util.GetLogger()

	p.shutdownFlag = true
	var reason string
	timeoutCh := make(chan string, 1)
	go func() {
		time.Sleep(10 * time.Second)
		log.Debug().Msg("Timeout reached")
		timeoutCh <- "TIMEOUT"
	}()

	shutdownDone := make(map[string]bool, len(p.channels))
	for k, ch := range p.channels {
		log.Debug().Msgf("Sending shutdown request to %s", k)
		ch <- "SHUTDOWN"
	}

	for {
		stop := false
		for k, ch := range p.channels {
			select {
			case <-ch:
				shutdownDone[k] = true
				if len(p.channels) == len(shutdownDone) {
					stop = true
					reason = "SHUTDOWN COMPLETE"
				}
			case <-timeoutCh:
				reason = "TIMEOUT"
				stop = true
			default:
			}
		}
		if stop {
			break
		}
	}

	return reason
}

func NewProcessManager(channel chan string) ProcessManager {
	var p ProcessManager
	p.inChannel = channel
	p.shutdownFlag = false
	p.channels = make(map[string]chan string, 10)

	return p
}
