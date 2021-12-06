package executer

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

// Create from jobdefinition. Filter out not needed for scheduling.
type jobTask struct {
	message.JobRequest
	Started      bool
	JobDef       objects.JobDefinition
	JobStepTasks []*jobStepTask // Write once read many
	JobStepState map[string]*jobStepStatus
	NextJobStep  string
	QueuedAt     int64
	RunAt        int64
	EndAt        int64
}

//
type jobStepTask struct {
	objects.JobStepDefinition
	PreSteps []string
}

type jobStepStatus struct {
	Running bool
	Done    bool
}
