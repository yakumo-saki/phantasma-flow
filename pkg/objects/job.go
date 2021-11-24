package objects

import "fmt"

//
type JobDefinition struct {
	ObjectBase
	Meta    ObjectMetaBase      `json:"meta"`
	JobMeta JobMetaInfo         `json:"jobmeta"`
	Steps   []JobStepDefinition `json:"steps"`
	Name    string              `json:"name"`
	Id      string              `json:"id"` // sha1sum Name. it is used for filename of job related.
}

type JobMetaInfo struct {
	Schedules []JobSchedule `json:"schedules"`
}

type JobSchedule struct {
	ScheduleType string `json:"type"` // cron
	Param        string `json:"param"`
}

type JobStepDefinition struct {
	Name     string `json:"name"`
	ExecType string `json:"execType"` // normal , internal
	Command  string `json:"command"`
	Script   string `json:"script"`
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
