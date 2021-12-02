package repository

func GetLogDirectory() string {
	return repo.paths.JobLog
}

func GetJobMetaDirectory() string {
	return repo.paths.JobMeta
}
