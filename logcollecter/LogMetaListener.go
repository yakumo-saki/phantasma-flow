package logcollecter

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/yakumo-saki/phantasma-flow/global"
	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type logMetaListenerParams struct {
	runId    string
	jobId    string
	alive    bool
	execChan chan message.ExecuterMsg
	ctx      context.Context
	cancel   context.CancelFunc
}

// This handles executer.ExecuterMsg
// Collect and save jobresult (executed job step result)
func (m *LogListenerModule) LogMetaListener(ctx context.Context) {
	NAME := "LogMetaListener"
	log := util.GetLoggerWithSource(m.GetName(), NAME)

	defer m.logChannelsWg.Done()

	loggerMap := make(map[string]*logMetaListenerParams) // runid -> loglistener
	waitGroup := sync.WaitGroup{}
	logInCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, NAME)

	for {
		select {
		case msg := <-logInCh:
			execMsg := msg.Body.(message.ExecuterMsg)

			listener, ok := loggerMap[execMsg.JobId] // JobIDで見ているのは、JobMeta fileがJobId単位だから
			if !ok || !listener.alive {
				log.Trace().Msgf("create meta listener for %s", execMsg.RunId)
				loglis := m.createJobLogMetaListenerParams(execMsg)
				loggerMap[execMsg.RunId] = loglis

				waitGroup.Add(1)
				go m.jobLogMetaListener(loglis, &waitGroup)
				listener = loglis
			}

			listener.execChan <- execMsg
		case <-ctx.Done():
			goto shutdown
		}
	}

shutdown:
	messagehub.Unsubscribe(messagehub.TOPIC_JOB_REPORT, NAME)
	for _, metalis := range loggerMap {
		metalis.cancel()
	}
	waitGroup.Wait()
	log.Info().Msgf("%s/%s stopped.", m.GetName(), NAME)
}

func (m *LogListenerModule) createJobLogMetaListenerParams(lm message.ExecuterMsg) *logMetaListenerParams {

	loglis := logMetaListenerParams{runId: lm.RunId, jobId: lm.JobId}
	ch := make(chan message.ExecuterMsg, 1)
	loglis.execChan = ch
	loglis.ctx, loglis.cancel = context.WithCancel(context.Background())
	return &loglis
}

// Trace single job
// ログファイルの競合を避けるために1 runId = 1goroutineとして走らせることにした
func (m *LogListenerModule) jobLogMetaListener(params *logMetaListenerParams, wg *sync.WaitGroup) {
	NAME := "jobLogMetaListener"
	log := util.GetLoggerWithSource(m.GetName(), NAME).
		With().Str("jobId", params.jobId).Logger()

	// 既存ログファイルオープン
	logDir := repository.GetJobMetaDirectory()
	filename := fmt.Sprintf("%s.yaml", time.Now().Format(global.DATETIME_FORMAT))
	logpath := path.Join(logDir, filename)

	// TODO load already existed
	var metaLog *logfile.JobMetaLog
	if true {
		metaLog = m.createEmptyJobLogMeta(params.jobId)
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
		case <-params.ctx.Done():
			log.Debug().Msg("Shutdown request received.")
			goto shutdown
		}
	}
shutdown:
	const RESULT_LOST = "job results are lost"
	params.alive = false

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

	wg.Done()
	log.Debug().Msgf("Stopped jobLogMetaListener for %s", params.jobId)
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
