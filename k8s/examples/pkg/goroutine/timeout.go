package goroutine

import (
	"fmt"
	"time"
)

// sT := 3 * time.Second
//
//	go FuncWithTimeout(func() {
//		time.Sleep(1 * time.Second)
//	}, sT)
func FuncWithTimeout(f func(), timeout time.Duration) {
	c := make(chan struct{}, 1)

	go func() {
		f()
		c <- struct{}{}
	}()
	select {
	case <-c:
		fmt.Println("finish success")
	case <-time.After(timeout):
		fmt.Println("time out")
	}
}
