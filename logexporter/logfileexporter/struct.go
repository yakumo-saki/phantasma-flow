package logfileexporter

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
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
	logChan  chan *objects.JobLogMessage
}
