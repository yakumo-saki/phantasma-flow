package executer

import "github.com/yakumo-saki/phantasma-flow/pkg/objects"

// Add new job
func (js *Executer) AddJob(jobDef objects.JobDefinition) {
	// j := job{}
	// j.id = jobDef.Id
	// j.jobMeta = jobDef.JobMeta
	// j.lastRun = 0
	// j.name = jobDef.Name

	js.mutex.Lock()
	defer js.mutex.Unlock()
	// js.jobs[j.id] = j
}
