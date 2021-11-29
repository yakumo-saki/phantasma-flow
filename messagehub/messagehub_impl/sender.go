package messagehub_impl

import (
	"context"
	"time"

	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Sender thread, single thread expected
func (hub *MessageHub) Sender(ctxptr *context.Context) {
	log := util.GetLoggerWithSource(hub.Name, "sender")
	defer hub.senderWaitGroup.Done()

	log.Debug().Msg("Sender started.")

	ctx := *ctxptr
	for {
		var msg *message.Message
		select {
		case <-ctx.Done():
			log.Info().Msg("Sender stopped.")
			return
		default:
			// wait for message
			c, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			m, err := hub.queue.DequeueOrWaitForNextElementContext(c)
			cancel()

			if err != nil {
				continue // nothing in queue
			}
			msg = m.(*message.Message)
		}

		topic := msg.Topic

		// cant use defer. this routine is not exit.
		hub.listenerMutex.Lock()

		liss, ok := hub.listeners.Load(topic)
		if !ok {
			log.Debug().Msgf("No listener for topic %s", topic)
			hub.listenerMutex.Unlock()
			continue
		}

		listerners := liss.(*[]listener)
		for _, lis := range *listerners {
			log.Trace().Str("topic", topic).Str("listener", lis.name).Msgf("%v", msg)
			lis.ch <- msg
		}
		hub.listenerMutex.Unlock()
	}
}
