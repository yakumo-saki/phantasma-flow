package logfileexporter

import (
	"context"
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/logcollecter/logfile"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// main method of LogListener
func (m *LogFileExporter) LogListener(ctx context.Context) {
	const NAME = "LogListener"
	log := util.GetLoggerWithSource(m.GetName(), NAME)
	defer m.logChannelsWg.Done()

	loggerMap := make(map[string]*logListenerParams) // runid -> loglistener
	waitGroup := sync.WaitGroup{}
	logCh := messagehub.Subscribe(messagehub.TOPIC_JOB_LOG, NAME)

	select {
	case msg := <-logCh:
		joblogMsg := msg.Body.(logfile.JobLogMessage)

		listener, ok := loggerMap[joblogMsg.RunId] // Job log fileはRunId単位
		if !ok || !listener.Alive {
			log.Trace().Msgf("create joblog listener for %s", joblogMsg.RunId)
			loglis := m.createJobLogListenerParams(joblogMsg)
			loggerMap[joblogMsg.RunId] = loglis

			waitGroup.Add(1)
			go m.joblogListener(loglis, &waitGroup)
			listener = loglis
		}

		listener.logChan <- joblogMsg

	case <-ctx.Done():
		goto shutdown

		// case <-messagehub.msg()
		// close it
	}

shutdown:
	log.Debug().Msg("Stopping all log listerners.")
	m.logCloseFunc.Range(func(key interface{}, cf interface{}) bool {
		log.Trace().Msgf("context stop %v", key)
		cancelFunc := cf.(context.CancelFunc)
		cancelFunc()
		return true
	})

	doneCh := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		time.Sleep(100 * time.Millisecond)
		waitGroup.Wait()
		close(ch)
	}(doneCh)

	select {
	case <-doneCh:
		log.Info().Msg("Stopping all log listerners completed")
	case <-time.After(10 * time.Second):
		log.Warn().Msg("Stopping all log listerners timeout")
	}

}

func (m *LogFileExporter) createJobLogListenerParams(lm logfile.JobLogMessage) *logListenerParams {

	loglis := logListenerParams{}
	loglis.RunId = lm.RunId
	loglis.JobId = lm.JobId
	ch := make(chan logfile.JobLogMessage, 1)
	loglis.logChan = ch
	loglis.Ctx, loglis.Cancel = context.WithCancel(context.Background())
	return &loglis
}
