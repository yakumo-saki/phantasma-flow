package metalistener

import (
	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

func (m *jobLogMetaListener) createEmptyJobLogMeta(jobId string) *logfile.JobMetaLog {

	jm := logfile.JobMetaLog{}
	jm.JobId = jobId
	jm.Kind = logfile.KIND_JOB_META
	jm.Meta = logfile.JobMetaMeta{}
	jm.Meta.NextJobNumber = 1
	jm.Results = []*logfile.JobMetaResult{}

	return &jm
}

func (m *jobLogMetaListener) createNewJobLogMetaResult(runId string, ver objects.ObjectVersion) *logfile.JobMetaResult {

	result := logfile.JobMetaResult{}
	result.JobNumber = -1 // invalid value.
	result.Success = false
	result.RunId = runId
	result.Version = ver
	result.StepResults = []*logfile.JobMetaStepResult{}

	return &result

}

func (m *jobLogMetaListener) createJobStepMetaResult(stepName string) *logfile.JobMetaStepResult {
	stepResult := logfile.JobMetaStepResult{}
	stepResult.StepName = stepName
	stepResult.ExitCode = -1
	stepResult.Success = false

	return &stepResult
}

// JobMetaResultは新しいもの順に記録したいので append slice newest first
func (m *jobLogMetaListener) appendMetaResult(results []*logfile.JobMetaResult,
	newResult *logfile.JobMetaResult) []*logfile.JobMetaResult {
	var slice []*logfile.JobMetaResult
	slice = append(slice, newResult)
	slice = append(slice, results...)

	return slice
}

func (m *jobLogMetaListener) findMetaResultByRunId(results []*logfile.JobMetaResult, runId string) *logfile.JobMetaResult {
	for _, jmr := range results {
		if jmr.RunId == runId {
			jmrCopy := jmr
			return jmrCopy
		}
	}

	return nil // not found
}

func (m *jobLogMetaListener) findStepResultByStepName(results []*logfile.JobMetaStepResult, stepName string) *logfile.JobMetaStepResult {
	for _, sr := range results {
		srCopy := sr
		if sr.StepName == stepName {
			return srCopy
		}
	}
	return nil // not found
}
