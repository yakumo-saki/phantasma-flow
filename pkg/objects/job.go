package objects

import "fmt"

//
type JobDefinition struct {
	ObjectBase
	Meta    ObjectMetaBase      `json:"meta"`
	JobMeta JobMetaInfo         `json:"jobmeta"`
	Steps   []JobStepDefinition `json:"steps"`
	Name    string              `json:"name"`
	Id      string              `json:"id"` // string-ID. it is used for filename of job related.
}

type JobMetaInfo struct {
	Schedules []JobSchedule `json:"schedules"`
}

type JobSchedule struct {
	ScheduleType string `json:"type"` // cron
	Param        string `json:"param"`
}

const JOB_EXEC_TYPE_COMMAND = "command"
const JOB_EXEC_TYPE_SCRIPT = "script"

type JobStepDefinition struct {
	Name        string `json:"name"`        // JobStep Name
	UseCapacity int    `json:"useCapacity"` // number how many capacity this step use. default 1
	ExecType    string `json:"execType"`    // JOB_EXEC_TYPE_*. internal, command, script
	Command     string `json:"command"`     // JOB_EXEC_TYPE_COMMAND only
	Script      string `json:"script"`      // JOB_EXEC_TYPE_SCRIPT only
}

func (nd JobDefinition) String() string {
	ret := fmt.Sprintf("Name: %s Meta: %v", nd.Name, nd.Meta)

	ret = ret + "\n"
	for _, st := range nd.Steps {
		ret = ret + fmt.Sprintf("Step: %v\n", st)
	}
	return ret
}

func (nd JobStepDefinition) String() string {
	return nd.Name
}
