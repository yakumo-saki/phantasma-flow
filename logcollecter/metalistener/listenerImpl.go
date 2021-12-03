package metalistener

import (
	"fmt"
	"io/ioutil"
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
func (m *MetaListener) jobLogMetaListener(params *logMetaListenerParams, wg *sync.WaitGroup) {
	NAME := "jobLogMetaListener"
	log := util.GetLoggerWithSource(m.GetName(), NAME).
		With().Str("jobId", params.JobId).Logger()

	defer wg.Done()

	// 既存ログファイルオープン
	logDir := repository.GetJobMetaDirectory()
	filename := fmt.Sprintf("%s.yaml", params.JobId)
	metaFilePath := path.Join(logDir, filename)

	var metaLog *logfile.JobMetaLog
	if util.IsFileExist(metaFilePath) {
		bytes, err := ioutil.ReadFile(metaFilePath)
		if err != nil {
			panic(err)
		}

		meta := &logfile.JobMetaLog{}
		err = yaml.Unmarshal(bytes, meta)
		if err == nil {
			metaLog = meta
		} else {
			log.Err(err).Msgf("JobMeta yaml is broken. Recreate %s", metaFilePath)
			metaLog = m.createEmptyJobLogMeta(params.JobId)
		}
	} else {
		metaLog = m.createEmptyJobLogMeta(params.JobId)
	}

	//
	findOrCreateMetaStep := func(jobResult *logfile.JobMetaResult, stepName string) (*logfile.JobMetaStepResult, bool) {
		// find for JobMetaResults
		sr := m.findStepResultByStepName(jobResult.Results, stepName)
		if sr != nil {
			return sr, false
		}

		new := m.createJobStepMetaResult(stepName)
		return new, true
	}

	for {
		var jobResult *logfile.JobMetaResult
		select {
		case msg, ok := <-params.execChan:
			if ok {
				l := log.With().Str("runId", msg.RunId).Logger()

				// find for JobMetaResults
				{
					jresult := m.findMetaResultByRunId(metaLog.Results, msg.RunId)
					if jresult != nil {
						jobResult = jresult
					} else {
						jobResult = m.createNewJobLogMetaResult(msg.RunId, msg.Version)
					}
				}

				switch msg.Reason {
				case message.JOB_START:
					jobResult = m.createNewJobLogMetaResult(msg.RunId, msg.Version)
					jobResult.JobNumber = metaLog.Meta.NextJobNumber
					jobResult.StartDateTime = util.GetDateTimeString()
					jobResult.EndDateTime = ""

					metaLog.Results = m.appendMetaResult(metaLog.Results, jobResult)
					metaLog.Meta.NextJobNumber++
				case message.JOB_END:
					jobResult.EndDateTime = util.GetDateTimeString()
				case message.JOB_STEP_START:
					stepResult := m.createJobStepMetaResult(msg.StepName)
					stepResult.StartDateTime = util.GetDateTimeString()
					jobResult.Results = append(jobResult.Results, *stepResult)
				case message.JOB_STEP_END:
					stepResult, created := findOrCreateMetaStep(jobResult, msg.StepName)
					if created {
						l.Warn().Str("stepName", msg.StepName).Msg("JOB_STEP_END received but JobMetaStepResult not found")
						jobResult.Results = append(jobResult.Results, *stepResult)
					}
					stepResult.EndDateTime = util.GetDateTimeString()
				}

				log.Debug().Msgf("%v", msg)
			} else {
				log.Debug().Msg("Shutdown request received via channel close")
				goto shutdown
			}

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
	f, err := os.OpenFile(metaFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Err(err).Msgf("failed to open job meta file %s. %s", metaFilePath, RESULT_LOST)
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
		log.Err(err).Msgf("failed to close %s.", metaFilePath)
	}

	log.Debug().Msgf("Stopped jobLogMetaListener for %s", params.JobId)
}

func (m *MetaListener) createEmptyJobLogMeta(jobId string) *logfile.JobMetaLog {

	jm := logfile.JobMetaLog{}
	jm.Kind = logfile.KIND_JOB_META
	jm.Meta = logfile.JobMetaMeta{}
	jm.Meta.NextJobNumber = 1
	jm.Results = []logfile.JobMetaResult{}

	return &jm
}

func (m *MetaListener) createNewJobLogMetaResult(runId string, ver objects.ObjectVersion) *logfile.JobMetaResult {

	result := logfile.JobMetaResult{}
	result.JobNumber = -1 // invalid value.
	result.Success = false
	result.RunId = runId
	result.Version = ver
	result.Results = []logfile.JobMetaStepResult{}

	return &result

}

func (m *MetaListener) createJobStepMetaResult(stepName string) *logfile.JobMetaStepResult {
	stepResult := logfile.JobMetaStepResult{}
	stepResult.StepName = stepName
	stepResult.ExitCode = -1
	stepResult.Success = false
	return &stepResult
}

// JobMetaResultは新しいもの順に記録したいので append slice newest first
func (m *MetaListener) appendMetaResult(results []logfile.JobMetaResult,
	newResult *logfile.JobMetaResult) []logfile.JobMetaResult {
	var slice []logfile.JobMetaResult
	slice = append(slice, *newResult)
	slice = append(slice, results...)

	return slice
}

func (m *MetaListener) findMetaResultByRunId(results []logfile.JobMetaResult, runId string) *logfile.JobMetaResult {
	for _, jmr := range results {
		if jmr.RunId == runId {
			return &jmr
		}
	}
	return nil // not found
}

func (m *MetaListener) findStepResultByStepName(results []logfile.JobMetaStepResult, stepName string) *logfile.JobMetaStepResult {
	for _, sr := range results {
		if sr.StepName == stepName {
			return &sr
		}
	}
	return nil // not found
}
