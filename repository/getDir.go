package repository

import "github.com/yakumo-saki/phantasma-flow/util"

// GetLogDirectory returns log directory.
// It is not change while running
func GetLogDirectory() string {
	checkRepoInitialized()
	return repo.paths.JobLog
}

// GetJobMetaDirectory returns job meta directory.
// It is not change while running
func GetJobMetaDirectory() string {
	checkRepoInitialized()
	return repo.paths.JobMeta
}

func checkRepoInitialized() {
	log := util.GetLoggerWithSource("repository")
	if repo == nil || !repo.Initialized {
		log.Error().Msg("Repository is not initialized!!")
		panic("Repository is not initialized!!")
	}
}
