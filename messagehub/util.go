package messagehub

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
)

// Wait for all message is sent.
// NOTE: This not blocks incoming new message.
func WaitForQueueEmpty(log *zerolog.Logger, hub *messagehub_impl.MessageHub) {
	for {
		if hub.GetQueueLength() == 0 {
			log.Debug().Msg("Wait for message hub done.")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
