package jobparser

import (
	"container/list"
	"errors"
	"fmt"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

var repo *repository.Repository

func BuildExecutableJob(jobId, runId string) (*list.List, error) {
	log := util.GetLoggerWithSource("jobParser", "BuildExecutableJob").
		With().Str("jobId", jobId).Logger()

	if repo == nil {
		repo = repository.GetRepository()
	}

	jobDef := repo.GetJobById(jobId)

	if jobDef == nil {
		log.Error().Msgf("No job found in repository.")
		return &list.List{}, errors.New("no job found in repository")
	}

	result, err := BuildFromJobDefinition(jobDef, jobId, runId)
	if err != nil {
		log.Err(err).Msgf("Failed to build.")
		return &list.List{}, err
	}

	return result, nil
}

// buildFromJobDefinition builds ExecutableJobs as list.List.
func BuildFromJobDefinition(jobDef *objects.JobDefinition, jobId, runId string) (*list.List, error) {
	if len(jobDef.Steps) == 0 {
		msg := fmt.Sprintf("JobSteps is empty, JobId=%s (%s)", jobDef.Id, jobDef.Name)
		return list.New(), errors.New(msg)
	}

	switch jobDef.JobMeta.ExecType {
	case objects.JOB_TYPE_PARA:
		panic("not implemented")
	case objects.JOB_TYPE_SEQ:
		return buildFromSequentialJobDef(jobDef, jobId, runId)
	case "":
		return buildFromSequentialJobDef(jobDef, jobId, runId)
	default:
		msg := fmt.Sprintf("Unknown jobMeta.execType %s, JobId=%s (%s)", jobDef.JobMeta.ExecType, jobDef.Id, jobDef.Name)
		panic(msg)
	}
}

func buildFromSequentialJobDef(jobDef *objects.JobDefinition, jobId, runId string) (*list.List, error) {
	result := list.New()
	var lastStep *ExecutableJobStep

	for idx, step := range jobDef.Steps {
		execStep := ExecutableJobStep{}
		util.DeepCopy(jobDef, &execStep)
		util.DeepCopy(&step, &execStep)
		util.DeepCopy(jobDef.Meta.Version, &execStep.Version)

		execStep.RunId = runId
		execStep.JobId = jobId
		setDefaultValues(idx, &execStep)

		// PreStep
		if lastStep == nil {
			// first step. no condition
			execStep.PreSteps = []string{}
		} else {
			execStep.PreSteps = []string{lastStep.Name}
		}

		result.PushBack(execStep)
		lastStep = &execStep
		fmt.Println(execStep)
	}

	return result, nil

}

func setDefaultValues(index int, execStep *ExecutableJobStep) {
	// stepname default=step{n} n = 1 ~
	execStep.Name = ifEmpty(execStep.Name, fmt.Sprintf("step%v", index+1))

	execStep.Node = ifEmpty(execStep.Name, "local")

	if execStep.UseCapacity == -1 {
		execStep.UseCapacity = 1
	}

}

func ifEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
