package messagehub_impl

import (
	"context"
	"sync"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
)

type listener struct {
	name string
	ch   chan *message.Message
}

type MessageHub struct {
	listeners       sync.Map
	listenerMutex   sync.Mutex // to read listeners. note: mutex can be per topic basis for performance
	queue           *goconcurrentqueue.FIFO
	senderStarted   bool
	senderWaitGroup sync.WaitGroup
	senderCtx       *context.Context
	senderCancel    *context.CancelFunc
	messageCount    uint64
	Name            string
}
