package messagehub_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
)

func TestUnsubscribe(t *testing.T) {
	count := 0

	hub := messagehub_impl.MessageHub{}
	hub.Initialize()
	messagehub.SetMessageHub(&hub)

	// add listener
	wg := sync.WaitGroup{}
	wg.Add(2)
	go Listen(&hub, &wg, "topic1", "listner1")
	go ListenerWithShutdownCh(&hub, &wg, "topic1", "shutdown", "listner2")
	wg.Wait()

	hub.StartSender()

	hub.Post("topic1", "post test")
	count++

	messagehub.WaitForQueueEmpty("wait 1")
	messagehub.Unsubscribe("topic1", "listner222") // this make warning log
	messagehub.Unsubscribe("topic1", "listner2")

	hub.Post("topic1", "post test. not for listener2")
	count++

	fmt.Printf("Total messages sent %d\n", count)

	hub.Post("topic1", MSG_EXIT)   // not send to listner2
	hub.Post("shutdown", MSG_EXIT) // only for listner2

	hub.Shutdown()

	assert := assert.New(t)
	assert.Equal(count, getCount("listner1"))
	assert.Equal(count-1, getCount("listner2"))
}

func ListenerWithShutdownCh(hub *messagehub_impl.MessageHub, wg *sync.WaitGroup, topic, shutdownTopic, myname string) {
	count := 0
	ch := hub.Subscribe(topic, myname)
	stopch := hub.Subscribe(shutdownTopic, myname)
	wg.Done()
	for {
		var v *message.Message
		select {
		case v = <-ch:
		case v = <-stopch:
		}
		fmt.Printf("[%s] received: %s\n", myname, v.Body)
		if v.Body == MSG_EXIT {
			fmt.Printf("%s received count %d\n", myname, count)
			close(ch)

			listenResult.Store(myname, count)
			return
		}
		count++
	}
}
