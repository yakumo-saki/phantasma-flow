package logcollecter

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// Create log listener thread and channel
// After created, Other modules can get chan by calling GetLogListener
// Call by jobscheduler
// TODO: place log of internal jobs (log cleaner etcetc) is normal directory?
// NOTE: create after shutdown is not care. because jobscheduler is stopped before logcollecter stops
func (m *LogListenerModule) CreateLogListener(runId string, jobId string) <-chan LogMessage {
	ret := make(chan LogMessage)
	m.logChannels.Store(runId, ret)

	ctx, cancelFunc := context.WithCancel(m.RootCtx)
	m.logCloseFunc.Store(runId, cancelFunc)

	m.logChannelsWg.Add(1)
	go m.logListener(ctx, runId, ret)

	return ret
}

// Get logging chan by runId
func (m *LogListenerModule) GetLogListener(runId string) (<-chan LogMessage, bool) {
	cha, found := m.logChannels.Load(runId)

	if !found {
		return nil, false
	} else {
		return cha.(chan LogMessage), true
	}
}

// Maybe collect and save jobresult (executed job step result)
func (m *LogListenerModule) logListener(ctx context.Context, runId string, logInCh <-chan LogMessage) {

	defer m.logChannelsWg.Done()

	// open log file if exist
	// datetimeStr := time.Now().Format("20060102150405")
	// filename := fmt.Sprintf("%s_%s_%s")
	// wait loop and receive -> append
	// close log file

	// return
}

func (m *LogListenerModule) logListenerCloser(ctx context.Context) {

	log := util.GetLoggerWithSource(m.GetName(), "LogListernerCloser")

	// TODO listen messagehub jobscheduling topic.
	// if end of job -> cancelFunc

	select {
	case <-ctx.Done():
		goto shutdown
		// case <-messagehub.msg()
		// close it
	}
shutdown:
	log.Debug().Msg("Stopping all log listerners.")
	m.logCloseFunc.Range(func(key interface{}, cf interface{}) bool {
		cancelFunc := cf.(context.CancelFunc)
		cancelFunc()
		return true
	})

	doneCh := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		m.logChannelsWg.Wait()
		close(ch)
	}(doneCh)

	select {
	case <-doneCh:
		log.Warn().Msg("Stopping all log listerners completed")
	case <-time.After(10 * time.Second):
		log.Warn().Msg("Stopping all log listerners timeout")
	}

}
