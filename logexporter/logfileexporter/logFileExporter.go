package logfileexporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/global"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type logFileExporter struct {
}

func (m *logFileExporter) GetName() string {
	return "logFileExporter"
}

// Collect and save jobresult (executed job step result)
// When start / stop, change params.Alive flag
// 考え方：
// * ログが送られてきたらRunIdごとに1 goroutine起動
// * 終了はある程度の時間ログが送られてこなければタイムアウトして自動終了
// * ジョブ終了をlistenすることもできるが、messagehubは順番が保証されないので
//   ジョブ終了→ログとかなると困るので終了はタイムアウトのみとする
func (m *logFileExporter) Start(params *logListenerParams, wg *sync.WaitGroup) {
	NAME := "main"
	log := util.GetLoggerWithSource(m.GetName(), NAME).With().
		Str("jobId", params.JobId).Str("runId", params.RunId).Logger()

	params.Alive = true

	defer wg.Done()

	needClose := true
	useEmergencyLog := false
	emLog := util.GetLoggerWithSource(m.GetName(), "emergency")

	// XXX need to find same runid logfile
	f, err := m.openLogFile(params.JobId, params.RunId)
	if err != nil {
		useEmergencyLog = true
		needClose = false
	}

	for {
		select {
		case <-time.After(60 * time.Second):
			log.Debug().Msg("Automatic close.")
			goto shutdown
		case msg, ok := <-params.logChan:
			if !ok {
				log.Debug().Msg("Channel closed.")
				goto shutdown
			}

			// log.Debug().Msgf("%v", msg)

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
			log.Err(err).Msgf("Logfile close error. %s", f.Name())
		}
	}

	params.Alive = false
	log.Debug().Msgf("%s stopped.", m.GetName())
}

func (m *logFileExporter) openLogFile(jobId, runId string) (*os.File, error) {
	existedFilename := m.findJobLogFile(runId)
	if existedFilename == "" {
		return m.createLogFile(jobId, runId)
	} else {
		// open existant log
		f, err := os.OpenFile(existedFilename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Err(err).Msgf("Open log file failed %s, output to server log", existedFilename)
		}
		return f, err
	}

}

func (m *logFileExporter) createLogFile(jobId, runId string) (*os.File, error) {

	datetimeStr := time.Now().Format(global.DATETIME_FORMAT)
	filename := fmt.Sprintf("%s_%s_%s.json", datetimeStr, runId, jobId)

	logDir := repository.GetLogDirectory()
	logfile := path.Join(logDir, filename)

	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Err(err).Msgf("Create log file failed %s, output to server log", logfile)
	}
	return f, err
}

func (m *logFileExporter) findJobLogFile(runId string) string {
	logDir := repository.GetLogDirectory()

	pattern := path.Join(logDir, fmt.Sprintf("*_%s_*.json", runId))
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		ptn := strings.Split(file, "_")
		if len(ptn) != 3 {
			continue
		}

		fmt.Println(ptn[0])

		return file
	}
	return ""
}
