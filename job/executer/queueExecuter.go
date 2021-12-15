package executer

import (
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

func (ex *Executer) queueExecuter(startWg, stopWg *sync.WaitGroup) {
	const NAME = "queueExecuter"
	log := util.GetLoggerWithSource(ex.GetName(), NAME)
	log.Info().Msgf("Starting %s/%s.", ex.GetName(), NAME)

	defer stopWg.Done()
	startWg.Done()

	for {
		select {
		case <-ex.RootCtx.Done():
			goto shutdown
		case <-time.After(1 * time.Second):
			ex.mutex.Lock()
			ex.mutex.Unlock()
		}
	}

shutdown:
	log.Debug().Msgf("%s/%s stopped.", ex.GetName(), NAME)
}
