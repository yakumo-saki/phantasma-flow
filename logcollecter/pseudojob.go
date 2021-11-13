package logcollecter

import (
	"fmt"
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// This is only for test purpose
func PseudoLogSender(shutdown <-chan string, stop <-chan string, out chan string) {
	log := util.GetLogger()
	log.Info().Msg("Pseudo log sender start")

	stopFlag := false
	for {

		select {
		case v := <-stop:
			log.Info().Msgf("STOP signal received %s", v)
			stopFlag = true
		case v := <-shutdown:
			log.Info().Msgf("SHUTDOWN signal received %s", v)
			stopFlag = true
		default:
			msg := "Msg: " + fmt.Sprint(time.Now().Unix())
			out <- msg
			time.Sleep(1 * time.Second)
			log.Debug().Msg("PSEUDO " + msg)
			log.Debug().Msg("next loop pseudo")
		}

		if stopFlag {
			break
		}
	}

	log.Info().Msg("Pseudo log sender stopped")
}
