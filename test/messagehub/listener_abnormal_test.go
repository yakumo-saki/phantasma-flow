package messagehub_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
)

// Test abnormal state of messagehub
func TestMessageHubAbnormal(t *testing.T) {
	count := 0

	hub := messagehub.MessageHub{}
	hub.Initialize()

	hub.StopSender() // Stop not started sender is ok

	msg := hub.NewMessage()
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

	assert := assert.New(t)
	assert.Equal(uint64(count), hub.GetMessageCount())
	assert.Equal(count-1, (getCount("listner1"))) // listener not count EXIT msg
}
