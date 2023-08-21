package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"example.com/aryaLee/k8s/pkg/common"
	"example.com/aryaLee/k8s/pkg/lock"
	"github.com/google/uuid"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)

	wg := sync.WaitGroup{}
	n := 3
	wg.Add(n)
	for i := 0; i < n; i++ {
		ID := uuid.New().String()
		go func() {
			defer wg.Done()
			runlockFunc(ID)
		}()
	}
	wg.Wait()
}

func runlockFunc(ID string) {
	// use a Go context so we can tell the leaderelection code when we
	// want to step down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	run := func(ctx context.Context) {
		// complete your controller loop here
		klog.Info("Controller loop...")

		n := rand.Intn(10)
		fmt.Println(fmt.Sprintf("ID: %s, currrent Time: %s, n: %d", ID, time.Now(), n))

		time.Sleep(time.Duration(n) * time.Second)
		fmt.Println(fmt.Sprintf("finish ID: %s, currrent Time: %s, n: %d", ID, time.Now(), n))
	}

	client := common.NewClient()
	l := lock.LeaseLock{
		Namespace: "vpc-system",
		Name:      "test",
		Client:    client,
		Func:      run,
		ID:        ID,
	}

	// listen for interrupts or the Linux SIGTERM signal and cancel
	// our context, which the leader election code will observe and
	// step down
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		klog.Info("Received termination, signaling shutdown")
		cancel()
	}()

	l.Run(ctx)
}
