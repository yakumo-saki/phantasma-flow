package objects

import "fmt"

//
type JobDefinition struct {
	ObjectBase `yaml:",inline"`
	Meta       ObjectMetaBase      `yaml:"meta"`
	JobMeta    JobMetaInfo         `yaml:"jobmeta"` // Job meta informations (schedules)
	Steps      []JobStepDefinition `yaml:"steps"`   // Jobstep definitions
	Name       string              `yaml:"name"`    // display name of this job
	Id         string              `yaml:"id"`      // string-ID. it is used for filename of job related.
}

const JOB_TYPE_SEQ = "sequential"
const JOB_TYPE_PARA = "parallel"

type JobMetaInfo struct {
	Schedules []JobSchedule `yaml:"schedules"` // Run schedules (empty ok)
	ExecType  string        `yaml:"execType"`  // Type of running jobsteps. JOB_TYPE_* default=sequential.
}

type JobSchedule struct {
	ScheduleType string `yaml:"type"` // cron
	Param        string `yaml:"param"`
}

const JOB_EXEC_TYPE_COMMAND = "command"
const JOB_EXEC_TYPE_SCRIPT = "script"

type JobStepDefinition struct {
	Name        string `yaml:"name"`                     // JobStep Name (optional on sequential job)
	UseCapacity int    `yaml:"useCapacity" default:"-1"` // number how many capacity this step use. default 1
	ExecType    string `yaml:"execType"`                 // JOB_EXEC_TYPE_*. internal, command, script, (optional when command or script is defined.)
	Command     string `yaml:"command"`                  // JOB_EXEC_TYPE_COMMAND only
	Script      string `yaml:"script"`                   // JOB_EXEC_TYPE_SCRIPT only
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
