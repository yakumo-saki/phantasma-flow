package repository

import "github.com/yakumo-saki/phantasma-flow/pkg/objects"

func (repo *Repository) GetAllJobDefs() []objects.JobDefinition {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	ret := []objects.JobDefinition{}
	copy(ret, repo.jobs)
	return ret
}
