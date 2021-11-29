package messagehub_impl

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// Wait for all message is sent.
// NOTE: This not blocks incoming new message.
func (hub *MessageHub) WaitForQueueEmpty(msg string) {
	log := util.GetLoggerWithSource(hub.Name, "waitForEmpty").
		With().Str("message", msg).Logger()
	for {
		if hub.GetQueueLength() == 0 {
			log.Debug().Msg("Wait for message hub done.")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
