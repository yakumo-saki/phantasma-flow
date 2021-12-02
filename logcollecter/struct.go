package logcollecter

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
)

type logListenerParamsBase struct {
	RunId  string
	JobId  string
	Alive  bool
	Ctx    context.Context
	Cancel context.CancelFunc
}

type logListenerParams struct {
	logListenerParamsBase
	logChan chan logfile.JobLogMessage // XXX: JOB Log
}

type logMetaListenerParams struct {
	logListenerParamsBase
	execChan chan message.ExecuterMsg
}
