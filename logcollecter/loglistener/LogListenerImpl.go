package loglistener

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Collect and save jobresult (executed job step result)
func (m *LogListener) joblogListener(params *logListenerParams, wg *sync.WaitGroup) {
	log := util.GetLoggerWithSource(m.GetName(), "logListener")

	defer wg.Done()

	datetimeStr := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s_%s", datetimeStr, params.RunId, params.JobId)

	logDir := repository.GetLogDirectory()
	logfile := path.Join(logDir, filename)

	needClose := true
	useEmergencyLog := false
	emLog := util.GetLoggerWithSource(m.GetName(), "emergency")

	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Err(err).Msgf("Open log file failed %s, output to server log", logfile)
		useEmergencyLog = true
		needClose = false
	}

	for {
		select {
		case msg := <-params.logChan:
			bytes, err := json.Marshal(msg)
			logmsg := ""
			if err != nil {
				log.Err(err).Msg("JSON Marshal error")
				logmsg = fmt.Sprint(msg) // fallback
			} else {
				logmsg = string(bytes)
			}

			if !useEmergencyLog {
				_, err = f.Write([]byte(logmsg + "\n"))
				if err != nil {
					log.Err(err).Msg("Log write error, use server log")
					useEmergencyLog = true
				}
			}

			// write log failed or open failed, then server log
			if useEmergencyLog {
				emLog.Info().Msg(logmsg)
			}
		case <-params.Ctx.Done():
			goto shutdown
		}
	}

shutdown:
	if needClose {
		if err := f.Close(); err != nil {
			log.Err(err).Msgf("Log close error %s", logfile)
		}
	}
}
