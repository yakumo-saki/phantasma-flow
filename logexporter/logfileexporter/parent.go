package logfileexporter

import (
	"context"
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// main method of LogListener
func (m *LogFileExporter) LogListener(ctx context.Context) {
	const NAME = "parent"
	log := util.GetLoggerWithSource(m.GetName(), NAME)
	defer m.logChannelsWg.Done()

	loggerMap := make(map[string]*logListenerParams) // runid -> loglistener
	exportWg := sync.WaitGroup{}
	logCh := messagehub.Subscribe(messagehub.TOPIC_JOB_LOG, NAME)

	for {
		select {
		case msg, ok := <-logCh:
			if !ok {
				goto shutdown // channel closed
			}

			joblogMsg := msg.Body.(*objects.JobLogMessage)

			listener, ok := loggerMap[joblogMsg.RunId] // Job log fileはRunId単位
			if !ok {
				log.Trace().Msgf("Create logFileExporter for runId %s", joblogMsg.RunId)
				loglis := m.createJobLogListenerParams(joblogMsg)
				loggerMap[joblogMsg.RunId] = loglis

				exportWg.Add(1)
				loglis.Alive = true
				go loglis.instance.Start(loglis, &exportWg)
				listener = loglis
			} else if !listener.Alive {
				listener.Alive = true
				exportWg.Add(1)
				go listener.instance.Start(listener, &exportWg)
			}
			listener.logChan <- joblogMsg
		case <-ctx.Done():
			goto shutdown
		}
	}
	// TODO clean up loggerMap every 30min #44

shutdown:
	log.Debug().Msg("Stopping all log listerners.")

	for id, loglis := range loggerMap {
		if loglis.Alive {
			log.Trace().Msgf("Stop %v", id)
			loglis.Cancel()
		} else {
			log.Trace().Msgf("Already stopped, Skip %v", id)
		}
	}

	doneCh := make(chan struct{}, 1)
	go func(ch chan struct{}, wg *sync.WaitGroup) {
		time.Sleep(100 * time.Millisecond)
		wg.Wait()
		close(ch)
	}(doneCh, &exportWg)

	select {
	case <-doneCh:
		log.Info().Msg("Stopping all log listerners completed")
	case <-time.After(10 * time.Second):
		log.Warn().Msg("Stopping all log listerners timeout")
	}

}

func (m *LogFileExporter) createJobLogListenerParams(lm *objects.JobLogMessage) *logListenerParams {

	loglis := logListenerParams{}
	loglis.RunId = lm.RunId
	loglis.JobId = lm.JobId
	ch := make(chan *objects.JobLogMessage, 1)
	loglis.logChan = ch
	loglis.Ctx, loglis.Cancel = context.WithCancel(context.Background())
	loglis.instance = logFileExporter{}
	return &loglis
}
