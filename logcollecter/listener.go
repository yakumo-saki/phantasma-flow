package logcollecter

import (
	"io"
	"net"
	"strings"

	"github.com/yakumo-saki/phantasma-flow/util"
)

func LogListener(conn net.Conn, shutdown <-chan string, stop chan string, logIn <-chan string) {
	log := util.GetLogger()

	defer conn.Close()
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
			log.Debug().Msg("Wait for channel")
			message, more := <-logIn
			if more {
				log.Debug().Str("message", message).Msg("msg received")
				_, err := io.Copy(conn, strings.NewReader(message+"\n"))
				if err != nil {
					log.Debug().Err(err).Msg("Send log failed or connection closed")
					stopFlag = true
				}
			} else {
				log.Debug().Msg("msg channel closed")
				break
			}
			log.Debug().Msg("next loop send_data")
		}

		if stopFlag {
			break
		}

	}

	stop <- "STOP"
	log.Info().Msg("send_data stopped")
}
