package messagehub

import (
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
	"github.com/yakumo-saki/phantasma-flow/pkg/messagehub"
)

var hub *messagehub_impl.MessageHub

func SetMessageHub(mhub *messagehub_impl.MessageHub) {
	hub = mhub
}

func Listen(topic string, name string) chan *messagehub.Message {
	return hub.Listen(topic, name)
}
func Post(topic string, body interface{}) {
	hub.Post(topic, body)
}
func NewMessage() *messagehub.Message {
	msg := messagehub.Message{}
	return &msg
}

func GetQueueLength() int {
	return hub.GetQueueLength()
}
