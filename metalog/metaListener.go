package metalog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	// "github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
	"gopkg.in/yaml.v2"
)

// Trace single Job ID
type jobLogMetaListener struct {
	Mutex           sync.Mutex
	MetaLogFilePath string
	MetaLog         *objects.JobMetaLog    // overall meta log data (=single yaml file)
	JobMetaLog      *objects.JobMetaResult // Single job run meta log (=single result)
}

func (m *jobLogMetaListener) GetName() string {
	return "jobLogMetaListener"
}

// Trace single Job ID
// データファイルの競合を避けるために1 JobID = 1goroutineとして走らせることにした
func (m *jobLogMetaListener) Start(params *logMetaListenerParams, wg *sync.WaitGroup) {
	NAME := "main"
	log := util.GetLoggerWithSource(m.GetName(), NAME).
		With().Str("jobId", params.JobId).Logger()

	defer wg.Done()

	// get config
	jobResultCount := 64 // default
	{
		bareLogcfg := repository.GetRepository().GetConfigByKind(objects.KIND_JOBLOG_CFG)
		if bareLogcfg != nil {
			logcfg := bareLogcfg.(objects.JoblogConfig)
			jobResultCount = logcfg.JobResultCount
		}
	}

	// 既存ログファイルオープン or 新規作成
	m.ReadOrCreateMetaLog(params.JobId)

	timeout := time.NewTimer(15 * time.Second)

	for {
		// clean stop timer
		if !timeout.Stop() {
			<-timeout.C
		}
		timeout.Reset(15 * time.Second)

		select {
		case <-params.Ctx.Done():
			log.Debug().Msg("Shutdown request received.")
			goto shutdown
		case <-timeout.C:
			log.Debug().Msg("Metalog timeout, automatic shutdown.")
			goto shutdown
		case msg, ok := <-params.execChan:
			if !ok {
				log.Debug().Msg("Shutdown request received via channel close")
				goto shutdown
			}

			m.Mutex.Lock() // For GetNextJobNumber()

			// find for JobMetaResults
			jresult := m.findMetaResultByRunId(m.MetaLog.Results, msg.RunId)
			if jresult != nil {
				m.JobMetaLog = jresult // found job result
			} else {
				m.JobMetaLog = m.createNewJobLogMetaResult(msg.RunId, msg.Version)
				m.MetaLog.Results = m.appendMetaResult(m.MetaLog.Results, m.JobMetaLog)
			}

			switch msg.Subject {
			case message.JOB_START:
				m.handleJobStart(msg)
			case message.JOB_STEP_START:
				m.handleJobStepStart(msg)
			case message.JOB_STEP_END:
				m.handleJobStepEnd(msg)
			case message.JOB_END:
				m.handleJobEnd(msg)
				log.Trace().Msg("Stop child process because job is ended.")

				// delete metalog results over limit
				if len(m.MetaLog.Results) > jobResultCount {
					m.MetaLog.Results = m.MetaLog.Results[:jobResultCount]
				}

				m.Mutex.Unlock()
				goto shutdown
			}

			m.Mutex.Unlock()
		}
	}
shutdown:
	params.Alive = false
	m.WriteMetaLogToFile()

	m.MetaLogFilePath = ""
	m.MetaLog = nil
	m.JobMetaLog = nil
	log.Debug().Msgf("Stopped %s for jobId %s", m.GetName(), params.JobId)
}

func (m *jobLogMetaListener) GetNextJobNumber() int {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if m.MetaLog == nil {
		panic("Request GetNextJobNumber but MetaLog is null")
	}
	return m.MetaLog.Meta.NextJobNumber
}

func (m *jobLogMetaListener) handleJobStart(msg *message.ExecuterMsg) {
	m.JobMetaLog.JobNumber = m.MetaLog.Meta.NextJobNumber
	m.JobMetaLog.StartDateTime = util.GetDateTimeString()
	m.JobMetaLog.EndDateTime = ""
	m.MetaLog.Meta.NextJobNumber++
}

// JOB_END
func (m *jobLogMetaListener) handleJobEnd(msg *message.ExecuterMsg) {
	m.JobMetaLog.EndDateTime = util.GetDateTimeString()
	m.JobMetaLog.Success = msg.Success
	m.JobMetaLog.Reason = msg.Reason
}

func (m *jobLogMetaListener) handleJobStepStart(msg *message.ExecuterMsg) {
	stepResult := m.createJobStepMetaResult(msg.StepName)
	stepResult.StartDateTime = util.GetDateTimeString()
	stepResults := append(m.JobMetaLog.StepResults, stepResult)
	m.JobMetaLog.StepResults = stepResults
}

// JOB_STEP_END
func (m *jobLogMetaListener) handleJobStepEnd(msg *message.ExecuterMsg) {
	l := log.With().Str("reason", msg.Subject).Str("runId", msg.RunId).Logger()

	stepResult := m.findStepResultByStepName(m.JobMetaLog.StepResults, msg.StepName)
	if stepResult == nil {
		stepResult = m.createJobStepMetaResult(msg.StepName)
		l.Warn().Str("stepName", msg.StepName).
			Msgf("JOB_STEP_END received but JobMetaStepResult not found")
		m.JobMetaLog.StepResults = append(m.JobMetaLog.StepResults, stepResult)

		for _, s := range m.JobMetaLog.StepResults {
			log.Trace().Msgf("found step %v", s)
		}
	}
	stepResult.EndDateTime = util.GetDateTimeString()
	stepResult.ExitCode = msg.ExitCode
	stepResult.Success = msg.Success
	stepResult.Reason = msg.Reason

}

// ReadOrCreateMetaLog reads existing yaml file or create new one
func (m *jobLogMetaListener) ReadOrCreateMetaLog(jobId string) {
	log := util.GetLoggerWithSource(m.GetName(), "ReadOrCreateMetaLog")

	logDir := repository.GetJobMetaDirectory()
	filename := fmt.Sprintf("%s.yaml", jobId)
	metaFilePath := path.Join(logDir, filename)
	m.MetaLogFilePath = metaFilePath

	meta, err := readMetaLogfile(jobId)
	if err != nil {
		log.Debug().Msgf("Create jobmeta file %s", metaFilePath)
		m.MetaLog = m.createEmptyJobLogMeta(jobId)
		m.WriteMetaLogToFile()
		return
	}

	m.MetaLog = meta
}

//
func readMetaLogfile(jobId string) (*objects.JobMetaLog, error) {
	log := util.GetLoggerWithSource("ReadMetaLogfile")

	logDir := repository.GetJobMetaDirectory()
	filename := fmt.Sprintf("%s.yaml", jobId)
	metaFilePath := path.Join(logDir, filename)

	// read existed yaml
	if !util.IsFileExist(metaFilePath) {
		return nil, errors.New("JobMeta yaml not found")
	}

	// unreadable yaml
	bytes, err := ioutil.ReadFile(metaFilePath)
	if err != nil {
		log.Err(err).Msgf("JobMeta yaml unreadable.")
		return nil, err
	}

	meta := &objects.JobMetaLog{}
	err = yaml.Unmarshal(bytes, meta)
	if err != nil {
		log.Err(err).Msgf("JobMeta yaml is broken.")
		return nil, err
	}

	// validation
	if meta.Kind != objects.KIND_JOB_META {
		log.Err(err).Msgf("JobMeta yaml has wrong kind %s.", meta.Kind)
		return nil, err
	}

	return meta, nil // success
}

func (m *jobLogMetaListener) WriteMetaLogToFile() {
	const RESULT_LOST = "job results are lost"

	metaFilePath := m.MetaLogFilePath
	// write yaml
	f, err := os.OpenFile(metaFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Err(err).Msgf("failed to open job meta file %s. %s", metaFilePath, RESULT_LOST)
		return
	}
	bytes, err := yaml.Marshal(m.MetaLog)
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
		log.Err(err).Msgf("failed to close %s.", metaFilePath)
	}

	// json, _ := json.MarshalIndent(m.MetaLog, "", "  ")
	// fmt.Println(string(json))
}
