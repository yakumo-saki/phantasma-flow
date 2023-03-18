package messagehub_test

import (
	"sync"
	"testing"

	"github.com/huandu/go-assert"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
)

// Test abnormal state of messagehub
func TestMessageHubAbnormal(t *testing.T) {
	a := assert.New(t)
	count := 0

	hub := messagehub_impl.MessageHub{}
	hub.Initialize()
	messagehub.SetMessageHub(&hub)

	hub.StopSender() // Stop not started sender is ok

	msg := messagehub.NewMessage()
	msg.Topic = "topic1"
	msg.Body = "ABC1"
	hub.PostMsg(msg)
	count++

	// add listener
	wg := sync.WaitGroup{}
	wg.Add(1)
	go Listen(&hub, &wg, "topic1", "listner1")
	wg.Wait()

	hub.Post("topic1", "post test")
	count++

	hub.StartSender()
	hub.StopSender() // pause sender

	hub.Post("topic1", "post message in paused state is OK")
	count++

	hub.StartSender()

	hub.Post("topic1", MSG_EXIT)
	count++
	hub.Shutdown()

	a.Equal(uint64(count), hub.GetMessageCount())
	a.Equal(count-1, (getCount("listner1"))) // listener not count EXIT msg
}
