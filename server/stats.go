package server

import (
	"sync/atomic"
)

func (sv *Server) GetConnectionCount() int32 {
	return atomic.LoadInt32(&sv.connections)
}
