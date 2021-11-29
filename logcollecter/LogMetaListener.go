package logcollecter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Collect and save jobresult (executed job step result)
func (m *LogListenerModule) jobMetaListener(ctx context.Context, runId string, jobId string, logInCh <-chan LogMessage) {
	NAME := "jobMetaListener"
	log := util.GetLoggerWithSource(m.GetName(), NAME)

	defer m.logChannelsWg.Done()

	messagehub.Listen(messagehub.TOPIC_JOB_REPORT, NAME)

	for {
		select {
		case msg := <-logInCh:
			bytes, err := json.Marshal(msg)
			logmsg := ""
			if err != nil {
				log.Err(err).Msg("JSON Marshal error")
				logmsg = fmt.Sprint(msg) // fallback
			} else {
				logmsg = string(bytes)
			}

		case <-ctx.Done():
			goto shutdown
		}
	}

shutdown:
}

// Check single job
func (m *LogListenerModule) jobMetaTracer(ctx context.Context, runId string, jobId string, logInCh <-chan LogMessage) {
	logDir := repository.GetJobMetaDirectory()

}
