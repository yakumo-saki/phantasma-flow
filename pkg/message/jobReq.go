package message

type JobRequest struct {
	Reason string // ADD , CHANGE , INITIAL
	JobId  string
	RunId  string
}
