package messagehub_test

import (
	"fmt"
	"testing"

	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
)

func TestShutdownWithoutAnything(t *testing.T) {
	hub := messagehub_impl.MessageHub{}
	hub.Initialize()
	hub.Shutdown()
	fmt.Println("TEST END")
}

func TestStopBeforeStart(t *testing.T) {
	hub := messagehub_impl.MessageHub{}
	hub.Initialize()
	hub.StopSender() // Stop not started sender is ok
	fmt.Println("TEST END")
}

func TestShutdownWhenEmpty(t *testing.T) {
	hub := messagehub_impl.MessageHub{}
	hub.Initialize()
	hub.StartSender()
	hub.Shutdown()
	fmt.Println("TEST END")
}
