package jobparser

import (
	"fmt"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

type ExecutableJobStep struct {
	objects.JobStepDefinition
	Version objects.ObjectVersion
	JobId   string
	RunId   string
	Node    string
}

func (step *ExecutableJobStep) GetId() string {
	return fmt.Sprintf("%s_%s_%s", step.JobId, step.RunId, step.Name)
}
