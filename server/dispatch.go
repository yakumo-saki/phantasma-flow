package server

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/pkg/server"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Connected socket handling thread
// move to module
func (sv *Server) dispatch(conn net.Conn) {
	log := util.GetLoggerWithSource(sv.GetName(), "dispatch")

	log.Debug().Msg("request_dispatcher")
	scanner := bufio.NewScanner(conn)
	logchannel := make(chan string, 100)
	stopChannel := make(chan string, 1)

	ctx, negotiationDone := context.WithCancel(context.Background())
	defer negotiationDone()

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			log.Error().Msg("Timeout until negotiation. Closing connection")
			conn.Close()
		}
	}(ctx)

	for scanner.Scan() {
		line := scanner.Text() // スキャンした内容を文字列で取得
		lineStr := strings.ToUpper(strings.TrimSpace(string(line)))

		log.Debug().Str("set-type", lineStr).Msg("Received from client")
		switch lineStr {
		case "LISTENER":
			negotiationDone()
			log.Debug().Msg("Start listener")
			go logcollecter.PseudoLogSender(nil, stopChannel, logchannel)
		case "COMMANDER":
			negotiationDone()
			log.Debug().Msg("Start commander")
		case "PING":
			negotiationDone()
			res := server.ResPong{}
			res.Message = "PONG"
			msgBytes, _ := json.Marshal(res)
			msg := string(msgBytes) + "\n" + server.MSG_SEPARATOR

			sentBytes, err := io.Copy(conn, bytes.NewBufferString(msg))
			if err != nil {
				log.Err(err).Msg("Send PONG response failed")
			}
			fmt.Println(bytes.NewBufferString(msg).Len())
			log.Debug().Int64("bytes", sentBytes).Msg("Sent")
			conn.Close()
			return
		default:
			negotiationDone()
			if strings.Contains(lineStr, "HTTP/1.") {
				msg := "HTTP/1.0 400 Bad Request\n\n" +
					"This is not HTTP(s) port."

				_, err := io.Copy(conn, bytes.NewBufferString(msg))
				if err != nil {
					log.Err(err).Msg("Send 'this is not http' response failed")
				}
				conn.Close()
			}
		}
	}
}
