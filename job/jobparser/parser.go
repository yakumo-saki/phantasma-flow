package jobparser

import (
	"container/list"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

var repo *repository.Repository

func BuildExecutableJob(jobId, runId string) (list.List, error) {
	log := util.GetLoggerWithSource("jobParser", "BuildExecutableJob").
		With().Str("jobId", jobId).Logger()

	if repo == nil {
		repo = repository.GetRepository()
	}

	jobDef := repo.GetJobById(jobId)

	if jobDef == nil {
		log.Error().Msgf("No job found in repository.")
		return list.List{}, errors.New("no job found in repository")
	}

	result, err := BuildFromJobDefinition(jobDef, jobId, runId)
	if err != nil {
		log.Err(err).Msgf("Failed to build.")
		return list.List{}, err
	}

	return result, nil
}

// buildFromJobDefinition builds ExecutableJobs as list.List.
func BuildFromJobDefinition(jobDef *objects.JobDefinition, jobId, runId string) (list.List, error) {
	result := list.List{}

	for idx, step := range jobDef.Steps {
		execStep := ExecutableJobStep{}
		err := copier.Copy(&execStep, &jobDef)
		if err != nil {
			return list.List{}, err
		}
		err = copier.Copy(&execStep, &step)
		if err != nil {
			return list.List{}, err
		}
		execStep.RunId = runId
		execStep.JobId = jobId
		setDefaultValues(idx, &execStep)

		result.PushBack(execStep)
	}

	return result, nil
}

func setDefaultValues(index int, execStep *ExecutableJobStep) {
	// stepname default=step{n} n = 1 ~
	execStep.Name = iif(execStep.Name, fmt.Sprintf("step%v", index+1))

	execStep.Node = iif(execStep.Name, "local")

}

func iif(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
