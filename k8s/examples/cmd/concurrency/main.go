package main

import (
	"time"

	"example.com/aryaLee/k8s/pkg/goroutine"
)

func main() {
	st := 2 * time.Second
	go goroutine.FuncWithTimeout(func() {
		time.Sleep(1 * time.Second)
	}, st)
	time.Sleep(st + time.Second)
}
