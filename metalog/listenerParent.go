package metalog

import (
	"context"
	"sync"
	"time"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// This handles executer.ExecuterMsg
// Collect and save jobresult (executed job step result)
func (m *MetaListener) LogMetaListener(ctx context.Context) {
	const NAME = "LogMetaListener"
	log := util.GetLoggerWithSource(m.GetName(), NAME)

	defer m.logChannelsWg.Done()

	loggerMap := make(map[string]*logMetaListenerParams) // JobId
	waitGroup := sync.WaitGroup{}
	jobRepoCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, NAME)

	for {
		select {
		case msg, msgok := <-jobRepoCh:
			if !msgok {
				goto shutdown
			}
			execMsg := msg.Body.(*message.ExecuterMsg)

			listener, ok := loggerMap[execMsg.JobId] // JobIDで見ているのは、JobMeta fileがJobId単位だから
			if !ok {
				log.Trace().Msgf("create meta listener for %s %s", execMsg.JobId, execMsg.RunId)
				logmetaParams := m.createJobLogMetaListenerParams(execMsg)
				logmetaParams.instance = jobLogMetaListener{}
				loggerMap[execMsg.JobId] = logmetaParams

				waitGroup.Add(1)
				go logmetaParams.instance.Start(logmetaParams, &waitGroup)
				logmetaParams.Alive = true

				listener = logmetaParams
			} else if !listener.Alive {
				log.Debug().Msgf("Restart meta listener for %s %s", execMsg.JobId, execMsg.RunId)
				waitGroup.Add(1)
				go listener.instance.Start(listener, &waitGroup)
				listener.Alive = true
			}

			// log.Trace().Msgf("MetaLog send to %s %s", execMsg.JobId, execMsg.RunId)
			listener.execChan <- execMsg
			// log.Trace().Msgf("MetaLog send OK to %s %s", execMsg.JobId, execMsg.RunId)
		case <-ctx.Done():
			goto shutdown
		}
	}

shutdown:
	messagehub.Unsubscribe(messagehub.TOPIC_JOB_REPORT, NAME)
	for _, metalis := range loggerMap {
		metalis.Cancel()
	}

	doneCh := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		waitGroup.Wait()
		close(ch)
	}(doneCh)

	select {
	case <-doneCh:
		log.Info().Msg("Stopping all jobLogMetaListeners completed")
	case <-time.After(10 * time.Second):
		log.Warn().Msg("Stopping all jobLogMetaListeners timeout")
	}

	log.Info().Msgf("%s/%s stopped.", m.GetName(), NAME)
}

func (m *MetaListener) createJobLogMetaListenerParams(lm *message.ExecuterMsg) *logMetaListenerParams {

	loglis := logMetaListenerParams{}
	loglis.RunId = lm.RunId
	loglis.JobId = lm.JobId
	loglis.Alive = false
	ch := make(chan *message.ExecuterMsg, 1)
	loglis.execChan = ch
	loglis.Ctx, loglis.Cancel = context.WithCancel(context.Background())
	return &loglis
}
