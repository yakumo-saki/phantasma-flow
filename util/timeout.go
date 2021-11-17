package util

import (
	"time"
)

// cha := make(chan string, 1)
// go Timeout(cha, 10) => after 10 sec. TIMEOUT has send
func Timeout(timeoutCh chan string, second time.Duration) {
	time.Sleep(second * time.Second)
	timeoutCh <- "TIMEOUT"
}
