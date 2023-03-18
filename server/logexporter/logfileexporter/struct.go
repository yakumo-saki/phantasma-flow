package logfileexporter

import (
	"github.com/yakumo-saki/phantasma-flow/logcollecter"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

type logListenerParams struct {
	logcollecter.LogCollecterParamsBase
	instance logFileExporter
	logChan  chan *objects.JobLogMessage
}
