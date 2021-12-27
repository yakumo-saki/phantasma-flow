package metalog

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

func (m *jobLogMetaListener) createEmptyJobLogMeta(jobId string) *objects.JobMetaLog {

	jm := objects.JobMetaLog{}
	jm.JobId = jobId
	jm.Kind = objects.KIND_JOB_META
	jm.Meta = objects.JobMetaMeta{}
	jm.Meta.NextJobNumber = 1
	jm.Results = []*objects.JobMetaResult{}

	return &jm
}

func (m *jobLogMetaListener) createNewJobLogMetaResult(runId string, ver objects.ObjectVersion) *objects.JobMetaResult {

	result := objects.JobMetaResult{}
	result.JobNumber = -1 // invalid value.
	result.Success = false
	result.RunId = runId
	result.Version = ver
	result.StepResults = []*objects.JobMetaStepResult{}

	return &result

}

func (m *jobLogMetaListener) createJobStepMetaResult(stepName string) *objects.JobMetaStepResult {
	stepResult := objects.JobMetaStepResult{}
	stepResult.StepName = stepName
	stepResult.ExitCode = -1
	stepResult.Success = false

	return &stepResult
}

// JobMetaResultは新しいもの順に記録したいので append slice newest first
func (m *jobLogMetaListener) appendMetaResult(results []*objects.JobMetaResult,
	newResult *objects.JobMetaResult) []*objects.JobMetaResult {
	var slice []*objects.JobMetaResult
	slice = append(slice, newResult)
	slice = append(slice, results...)

	return slice
}

func (m *jobLogMetaListener) findMetaResultByRunId(results []*objects.JobMetaResult, runId string) *objects.JobMetaResult {
	for _, jmr := range results {
		if jmr.RunId == runId {
			jmrCopy := jmr
			return jmrCopy
		}
	}

	return nil // not found
}

func (m *jobLogMetaListener) findStepResultByStepName(results []*objects.JobMetaStepResult, stepName string) *objects.JobMetaStepResult {
	for _, sr := range results {
		srCopy := sr
		if sr.StepName == stepName {
			return srCopy
		}
	}
	return nil // not found
}
