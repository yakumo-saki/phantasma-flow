package objects

// internal
type JobState struct {
	JobId      string
	RunId      string
	Name       string
	JobState   string
	StepStates []JobStepStates
}

type JobStepStates struct {
	Name  string
	State string
}
