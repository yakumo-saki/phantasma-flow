package messagehub_bench

import (
	"fmt"
	"sync"
	"testing"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func TestBenchMessageHubSync(t *testing.T) {
	const LISTENER_MAX = 10

	// null sending
	hub := messagehub.MessageHub{}
	hub.Initialize()

	wg := sync.WaitGroup{}
	wg.Add(LISTENER_MAX)
	for i := 0; i < LISTENER_MAX; i++ {
		go ListenBench(&hub, &wg, "topic1", fmt.Sprintf("listner%d", i))
	}
	wg.Wait()

	hub.StartSender()

	timeoutCh := make(chan string, 1)
	go util.Timeout(timeoutCh, 5)

	fmt.Println("Start benchmark for sync post ")

	stop := false
	count := 0
	for {
		select {
		case <-timeoutCh:
			stop = true
		default:
			count++
			hub.Post("topic1", fmt.Sprintf("message_%d", count))
			hub.GetMessageCount()
			hub.GetQueueLength()
		}

		if stop {
			break
		}
	}

	hub.StopSender()
	fmt.Println("count for sync post ", count)
	fmt.Println("count for sent messages ", count-hub.GetQueueLength())

	hub.Post("topic1", MSG_EXIT)
}
