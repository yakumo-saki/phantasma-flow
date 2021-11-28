package messagehub

import (
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
)

var hub *messagehub_impl.MessageHub

func SetMessageHub(mhub *messagehub_impl.MessageHub) {
	hub = mhub
}

func Listen(topic string, name string) chan *message.Message {
	return hub.Listen(topic, name)
}
func Post(topic string, body interface{}) {
	hub.Post(topic, body)
}
func NewMessage() *message.Message {
	msg := message.Message{}
	return &msg
}

func GetQueueLength() int {
	return hub.GetQueueLength()
}
