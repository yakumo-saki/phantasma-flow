package objects

import "fmt"

//
type JobDefinition struct {
	ObjectBase
	Meta    ObjectMetaBase      `json:"meta"`
	JobMeta JobMetaInfo         `json:"jobmeta"` // Job meta informations (schedules)
	Steps   []JobStepDefinition `json:"steps"`   // Jobstep definitions
	Name    string              `json:"name"`    // display name of this job
	Id      string              `json:"id"`      // string-ID. it is used for filename of job related.
}

const JOB_TYPE_SEQ = "sequential"
const JOB_TYPE_PARA = "parallel"

type JobMetaInfo struct {
	Schedules []JobSchedule `json:"schedules"` // Run schedules (empty ok)
	ExecType  string        `json:"execType"`  // Type of running jobsteps. JOB_TYPE_* default=sequential.
}

type JobSchedule struct {
	ScheduleType string `json:"type"` // cron
	Param        string `json:"param"`
}

const JOB_EXEC_TYPE_COMMAND = "command"
const JOB_EXEC_TYPE_SCRIPT = "script"

type JobStepDefinition struct {
	Name        string `json:"name"`                     // JobStep Name (optional on sequential job)
	UseCapacity int    `json:"useCapacity" default:"-1"` // number how many capacity this step use. default 1
	ExecType    string `json:"execType"`                 // JOB_EXEC_TYPE_*. internal, command, script, (optional when command or script is defined.)
	Command     string `json:"command"`                  // JOB_EXEC_TYPE_COMMAND only
	Script      string `json:"script"`                   // JOB_EXEC_TYPE_SCRIPT only
}

func (jd JobDefinition) String() string {
	ret := fmt.Sprintf("Name: %s Meta: %v", jd.Name, jd.Meta)

	ret = ret + "\n"
	for _, st := range jd.Steps {
		ret = ret + fmt.Sprintf("Step: %v\n", st)
	}
	return ret
}

func (st JobStepDefinition) String() string {
	ret := fmt.Sprintf("StepName: %s Cap: %v ExeType: %s Cmd/Script: %s%s",
		st.Name, st.UseCapacity, st.ExecType, st.Command, st.Script)
	return ret
}
