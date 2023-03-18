package metalog

// This file' objective
// Provide direct access to nodeManager instance only for JobScheduler.

var logMetaManager *LogMetaManager

// GetInstance returns LogMetaManager instance
func GetInstance() *LogMetaManager {

	if logMetaManager == nil {
		logMetaManager = &LogMetaManager{}
		return logMetaManager
	}

	return logMetaManager
}
