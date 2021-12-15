package executer

var executer *Executer

func GetInstance() *Executer {
	if executer == nil {
		executer = &Executer{}
	}

	return executer
}

// GetJobQueueLength returns running jobs.
// note this cause mutex lock.
func GetJobQueueLength() int {
	return GetInstance().getJobQueueLength()
}
