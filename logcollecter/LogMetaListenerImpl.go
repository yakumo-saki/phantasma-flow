package logcollecter

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Trace single job
// ログファイルの競合を避けるために1 JobID = 1goroutineとして走らせることにした
func (m *LogListenerModule) jobLogMetaListener(params *logMetaListenerParams, wg *sync.WaitGroup) {
	NAME := "jobLogMetaListener"
	log := util.GetLoggerWithSource(m.GetName(), NAME).
		With().Str("jobId", params.JobId).Logger()

	defer wg.Done()

	// 既存ログファイルオープン
	logDir := repository.GetJobMetaDirectory()
	filename := fmt.Sprintf("%s.yaml", params.JobId)
	logpath := path.Join(logDir, filename)

	// TODO load already existed
	var metaLog *logfile.JobMetaLog
	if true {
		metaLog = m.createEmptyJobLogMeta(params.JobId)
	}

	for {
		select {
		case msg := <-params.execChan:
			switch msg.Reason {
			case message.JOB_START:
				jobResult := m.createNewJobLogMetaResult(msg.RunId, msg.Version)
				jobResult.JobNumber = metaLog.Meta.NextJobNumber
				metaLog.Results = append(metaLog.Results, *jobResult)
				metaLog.Meta.NextJobNumber++
			}

			log.Debug().Msgf("%v", msg)

		case <-time.After(15 * time.Second):
			log.Debug().Msg("timeout, auto close")
			goto shutdown
		case <-params.Ctx.Done():
			log.Debug().Msg("Shutdown request received.")
			goto shutdown
		}
	}
shutdown:
	const RESULT_LOST = "job results are lost"
	params.Alive = false

	// write yaml
	f, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Err(err).Msgf("failed to open job meta file %s. %s", logpath, RESULT_LOST)
		return
	}
	bytes, err := yaml.Marshal(metaLog)
	if err != nil {
		log.Err(err).Msgf("failed to serialize job meta result. %s", RESULT_LOST)
		return
	}
	_, err = f.Write(bytes)
	if err != nil {
		log.Err(err).Msgf("failed to Write job meta result. %s", RESULT_LOST)
	}
	err = f.Close()
	if err != nil {
		log.Err(err).Msgf("failed to close %s.", logpath)
	}

	log.Debug().Msgf("Stopped jobLogMetaListener for %s", params.JobId)
}

func (m *LogListenerModule) createEmptyJobLogMeta(jobId string) *logfile.JobMetaLog {

	jm := logfile.JobMetaLog{}
	jm.Kind = logfile.KIND_JOB_META
	jm.Meta = logfile.JobMetaMeta{}
	jm.Meta.NextJobNumber = 1
	jm.Results = []logfile.JobMetaResult{}

	return &jm
}

func (m *LogListenerModule) createNewJobLogMetaResult(runId string, ver objects.ObjectVersion) *logfile.JobMetaResult {

	result := logfile.JobMetaResult{}
	result.JobNumber = -1 // invalid value.
	result.Success = false
	result.RunId = runId
	result.Version = ver
	result.Results = []logfile.JobMetaStepResult{}

	return &result

}
