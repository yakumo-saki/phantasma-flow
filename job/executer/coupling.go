package executer

var executer *Executer

// GetInstance returns instance of Executer.
// Do not use executer.Executer{} because executer should be singleton
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
