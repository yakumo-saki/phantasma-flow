package messagehub_impl

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/yakumo-saki/phantasma-flow/pkg/messagehubObjects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type listener struct {
	name string
	ch   chan *messagehubObjects.Message
}

type MessageHub struct {
	listeners       sync.Map
	listenerMutex   sync.Mutex // to read listeners TODO: mutex can be per topic basis for performance
	queue           *goconcurrentqueue.FIFO
	senderCtx       *context.Context
	senderCancel    *context.CancelFunc
	senderWaitGroup sync.WaitGroup
	messageCount    uint64
}

func (hub *MessageHub) Initialize() {
	hub.listeners = sync.Map{}
	hub.queue = goconcurrentqueue.NewFIFO()
	hub.senderWaitGroup = sync.WaitGroup{}

	go hub.reportQueueLength()
}

func (hub *MessageHub) reportQueueLength() int {
	log := util.GetLogger()
	for {
		time.Sleep(30 * time.Second)
		log.Info().Msgf("Queue length: %d", hub.queue.GetLen())
	}
}

func (hub *MessageHub) StartSender() {
	if hub.senderCtx == nil {
		senderCtx, cancel := context.WithCancel(context.Background())
		hub.senderCtx = &senderCtx
		hub.senderCancel = &cancel
	}

	hub.senderWaitGroup.Add(1)
	go hub.Sender(hub.senderCtx)
}

// Stop sender thread. (Not waiting all queue done)
func (hub *MessageHub) StopSender() {
	log := util.GetLogger()
	if hub.senderCtx == nil { // not start senders and shutdown
		log.Info().Msgf("StopSender: No senders started. Nothing to do.")
		return
	}

	log.Debug().Msgf("Wait for stopping all senders.")

	cancel := *hub.senderCancel
	cancel()

	hub.senderWaitGroup.Wait()
	log.Info().Msgf("Shutdown all senders done.")
	hub.senderCancel = nil
	hub.senderCtx = nil
}

// Block new post and wait for queue become empty
func (hub *MessageHub) Shutdown() {
	log := util.GetLogger()

	// Immediate shutdown, when called shutdown in sender stopped state
	if hub.senderCtx == nil {
		log.Warn().Msgf("Shutdown immediate. because of no sender started. Wait queue left: %d", hub.queue.GetLen())
		return
	}

	// context to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// wait for queue flushed
	stop := false
	for {
		if hub.queue.GetLen() == 0 {
			log.Info().Msgf("Empty wait queue. continue shutdown")
			stop = true
			break
		}

		select {
		case <-ctx.Done():
			log.Warn().Msgf("Shutdown timeout. force shutdown. Wait queue left: %d", hub.queue.GetLen())
			stop = true
		case <-time.After(3 * time.Second):
			left := hub.queue.GetLen()
			log.Info().Msgf("Shutdown in progress. Wait queue left: %d", left)
			stop = (left == 0)
		}

		if stop {
			break
		}
	}

	hub.StopSender()
}

// go Sender()
func (hub *MessageHub) Sender(ctxptr *context.Context) {
	log := util.GetLogger()
	defer hub.senderWaitGroup.Done()

	log.Debug().Msg("Sender started.")

	ctx := *ctxptr
	for {
		var msg *messagehubObjects.Message
		select {
		case <-ctx.Done():
			log.Info().Msg("Sender stopped.")
			return
		default:
			c, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			m, err := hub.queue.DequeueOrWaitForNextElementContext(c)
			cancel()

			if err != nil {
				continue
			}
			msg = m.(*messagehubObjects.Message)
		}

		topic := msg.Topic

		// cant use defer. this routine is not exit.
		hub.listenerMutex.Lock()

		liss, ok := hub.listeners.Load(topic)
		if !ok {
			log.Debug().Msgf("No listener for topic %s", topic)
			hub.listenerMutex.Unlock()
			continue
		}

		listerners := liss.(*[]listener)
		for _, lis := range *listerners {
			lis.ch <- msg
		}
		hub.listenerMutex.Unlock()
	}
}

func (hub *MessageHub) GetQueueLength() int {
	return hub.queue.GetLen()
}

func (hub *MessageHub) GetMessageCount() uint64 {
	return atomic.LoadUint64(&hub.messageCount)
}

// Post() is add message to queue. no need to call as goroutine
// post(msg) is available.
// async / sync is up to you.
func (hub *MessageHub) Post(topic string, body interface{}) {
	hub.PostMsg(&messagehubObjects.Message{Topic: topic, Body: body})
}

func (hub *MessageHub) PostMsg(msg *messagehubObjects.Message) {
	hub.queue.Enqueue(msg)
	atomic.AddUint64(&hub.messageCount, 1)
}

func (hub *MessageHub) Listen(topic string, name string) chan *messagehubObjects.Message {
	log := util.GetLogger()

	hub.listenerMutex.Lock()
	defer hub.listenerMutex.Unlock()

	arr, ok := hub.listeners.Load(topic)
	array := &[]listener{}
	if ok {
		array = arr.(*[]listener)
	}

	ch := make(chan *messagehubObjects.Message, 1)
	newListener := listener{name: name, ch: ch}
	ls := append(*array, newListener)

	hub.listeners.Store(topic, &ls)

	log.Debug().Msgf("New listener added name:%s topic:%s listeners:%d",
		name, topic, len(ls))

	return ch
}
