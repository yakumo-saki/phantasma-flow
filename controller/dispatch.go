package controller

import (
	"bufio"
	"net"
	"strings"
	"time"

	job "github.com/yakumo-saki/phantasma-flow/jobscheduler"
	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func Dispatch(conn net.Conn, shutdownChannel <-chan string) {
	log := util.GetLogger()

	log.Debug().Msg("request_dispatcher")
	scanner := bufio.NewScanner(conn)
	logchannel := make(chan string, 100)
	stopChannel := make(chan string, 1)

	start := time.Now()
	for scanner.Scan() {
		line := scanner.Text() // スキャンした内容を文字列で取得
		lineStr := strings.ToUpper(strings.TrimSpace(string(line)))

		log.Debug().Str("set-type", lineStr).Msg("Received")
		if lineStr == "LISTENER" {
			log.Debug().Msg("Start listener")
			go logcollecter.LogListener(conn, shutdownChannel, stopChannel, logchannel)
			go logcollecter.PseudoLogSender(shutdownChannel, stopChannel, logchannel)
		} else if lineStr == "COMMANDER" {
			log.Debug().Msg("Start commander")
			go job.RequestHandler(conn, shutdownChannel, stopChannel, logchannel)
		}

		if time.Since(start).Seconds() > 10 {
			log.Error().Msg("Timeout waiting first message. Closing connection")
			conn.Close()
			break
		}
	}
}
