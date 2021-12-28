package message

import "fmt"

type Message struct {
	Topic string
	Body  interface{}
}

func (m *Message) String() string {
	msg := fmt.Sprintf("Msg: Topic=%s Body=%v", m.Topic, m.Body)
	return msg
}
