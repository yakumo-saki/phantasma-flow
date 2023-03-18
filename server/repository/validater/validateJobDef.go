package validater

import (
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
)

func ValidateAllJobDef() error {

	jobdefs := repository.GetRepository().GetAllJobDefs()

	for _, jd := range jobdefs {
		err := ValidateJobDef(jd)
		if err != nil {
			panic(err.Error())
		}

	}

	return nil
}

func ValidateJobDef(jobDef objects.JobDefinition) error {
	_, err := jobparser.BuildFromJobDefinition(&jobDef, "dummy")
	if err != nil {
		return err
	}
	return nil
}
