package messagehub_bench

import (
	"fmt"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
)

const MSG_EXIT = "exit"

func ListenBench(hub *messagehub_impl.MessageHub, wg *sync.WaitGroup, topic, myname string) {
	count := 0
	ch := hub.Listen(topic, myname)
	wg.Done()
	fmt.Printf("%s listen %s ok\n", myname, topic)
	for {
		v := <-ch
		count++
		if v.Body == MSG_EXIT {
			fmt.Printf("%s received count %d\n", myname, count)
			return
		}
	}
}
