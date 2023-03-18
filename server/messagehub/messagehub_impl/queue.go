package messagehub_impl

import (
	"sync/atomic"

	"github.com/yakumo-saki/phantasma-flow/pkg/message"
)

func (hub *MessageHub) GetQueueLength() int {
	return hub.queue.GetLen()
}

func (hub *MessageHub) GetMessageCount() uint64 {
	return atomic.LoadUint64(&hub.messageCount)
}

// Post() is add message to queue. no need to call as goroutine
// post(msg) is available.
func (hub *MessageHub) Post(topic string, body interface{}) {
	hub.PostMsg(&message.Message{Topic: topic, Body: body})
}

func (hub *MessageHub) PostMsg(msg *message.Message) {
	hub.queue.Enqueue(msg)
	atomic.AddUint64(&hub.messageCount, 1)
}
