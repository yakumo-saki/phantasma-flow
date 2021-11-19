package server

// Transport = JSON
// Req = client -> server
// Res = server -> client

const JOBID_ALL = "<ALL>"
const MSG_SEPARATOR = "\x00\x00\x00\n"

type AuthInfo struct {
	Token string `json:"token"`
}

type ReqSelectMode struct {
	AuthInfo
	Mode  string `json:"mode"`
	JobId string `json:"jobId"` // Job Log mode Only. "<ALL>" for all log.
}

type ResPong struct {
	Message string `json:"message"` // PONG
}
