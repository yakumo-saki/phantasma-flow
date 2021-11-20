package messagehub_test

import (
	"fmt"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
)

const MSG_EXIT = "exit"

var listenResult sync.Map

func Listen(hub *messagehub_impl.MessageHub, wg *sync.WaitGroup, topic, myname string) {
	count := 0
	ch := hub.Listen(topic, myname)
	wg.Done()
	for {
		v := <-ch
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

func ListenBench(hub *messagehub_impl.MessageHub, topic, myname string) {
	count := 0
	ch := hub.Listen(topic, myname)
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

func getCount(name string) int {
	r, ok := listenResult.Load(name)
	if !ok {
		return -1
	}
	return r.(int)
}
