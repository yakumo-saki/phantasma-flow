package jobparser

import (
	"fmt"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

type ExecutableJobStep struct {
	objects.JobStepDefinition
	Version  objects.ObjectVersion
	JobId    string
	RunId    string
	Node     string   // node id of running.
	PreSteps []string // Stepnames they must be completed before this job step run.
}

func (step ExecutableJobStep) String() string {
	msg := fmt.Sprintf("JobId:%s(v%v.%v) RunId:%s StepName:%s Cap: %v CMD:'%s' SCRIPT:'%s'",
		step.JobId, step.Version.Major, step.Version.Minor,
		step.RunId, step.Name,
		step.UseCapacity,
		step.Command, step.Script)
	return msg
}

func (step *ExecutableJobStep) GetId() string {
	return fmt.Sprintf("%s_%s_%s", step.JobId, step.RunId, step.Name)
}
