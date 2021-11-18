package messagehub_test

import (
	"fmt"
	"testing"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func TestBenchMessageHubSync(t *testing.T) {
	// null sending
	hub := messagehub.MessageHub{}
	hub.Initialize()
	// hub.StartSender()   // stop due to No listener for topic "topicname" log is weird

	timeoutCh := make(chan string, 1)
	go util.Timeout(timeoutCh, 1)

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

	fmt.Println("count for sync post ", count)
	fmt.Println("count for sent messages ", count-hub.GetQueueLength())

}
