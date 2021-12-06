package logfileexporter

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
)

type LogCollecterParamsBase struct {
	RunId  string
	JobId  string
	Alive  bool
	Ctx    context.Context
	Cancel context.CancelFunc
}

type logListenerParams struct {
	logcollecter.LogCollecterParamsBase
	instance logFileExporter
	logChan  chan logfile.JobLogMessage // XXX: JOB Log
}
