package metalistener

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
			log.Debug().Msgf("received %v", msg)
			execMsg := msg.Body.(message.ExecuterMsg)

			listener, ok := loggerMap[execMsg.JobId] // JobIDで見ているのは、JobMeta fileがJobId単位だから
			if !ok {
				log.Trace().Msgf("create meta listener for %s %s", execMsg.JobId, execMsg.RunId)
				loglis := m.createJobLogMetaListenerParams(execMsg)
				loggerMap[execMsg.JobId] = loglis

				waitGroup.Add(1)
				go m.jobLogMetaListener(loglis, &waitGroup)
				listener = loglis
			} else if !listener.Alive {
				log.Trace().Msgf("Restart meta listener for %s %s", execMsg.JobId, execMsg.RunId)
				loglis := loggerMap[execMsg.JobId]
				waitGroup.Add(1)
				go m.jobLogMetaListener(loglis, &waitGroup)
			}

			listener.execChan <- execMsg
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

func (m *MetaListener) createJobLogMetaListenerParams(lm message.ExecuterMsg) *logMetaListenerParams {

	loglis := logMetaListenerParams{}
	loglis.RunId = lm.RunId
	loglis.JobId = lm.JobId
	loglis.Alive = false
	ch := make(chan message.ExecuterMsg, 1)
	loglis.execChan = ch
	loglis.Ctx, loglis.Cancel = context.WithCancel(context.Background())
	return &loglis
}
