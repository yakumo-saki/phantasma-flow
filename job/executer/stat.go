package executer

func (ex *Executer) getJobQueueLength() int {
	ex.mutex.Lock()
	defer ex.mutex.Unlock()

	return len(ex.jobQueue)
}
