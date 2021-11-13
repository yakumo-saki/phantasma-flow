package jobscheduler

import (
	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type jobScheduler struct {
	globalCh chan string
	stopCh   <-chan string
}

var scheduler jobScheduler

func (js *jobScheduler) RequestHandler() {
	log := util.GetLogger()
	log.Debug().Msg("JobScheduler start")
}

func (js *jobScheduler) Loop() {
	stopFlag := false

	for {
		select {
		case v := <-js.stopCh:
			log.Info().Msgf("STOP signal received %s", v)
			stopFlag = true
		case v := <-js.globalCh:
			log.Info().Msgf("SHUTDOWN signal received %s", v)
			stopFlag = true
		default:
		}

		if stopFlag {
			break
		}
	}
	log.Info().Msg("JobScheduler stopped")
}

func (js *jobScheduler) Start() {
	log := util.GetLogger()

	log.Info().Msg("Starting JobScheduler.")
	go js.Loop()
	log.Info().Msg("Started JobScheduler.")
}

func Initialize() {
	repository.GetConfig()

}

func Start(globalCh chan string, stop <-chan string, out chan string) {
	scheduler.globalCh = globalCh
	scheduler.stopCh = stop
}
