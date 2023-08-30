package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	cancelCtx(ctx)
	fmt.Println("ctx error", ctx.Err())
	fmt.Println("ctx cause", context.Cause(ctx))
}

func cancelCtx(ctx context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	go sleepFunc(ctx, 1)
	time.Sleep(3 * time.Second)
	cancel(errors.New("ctx cancel test"))
	time.Sleep(5 * time.Second)
}

func sleepFunc(ctx context.Context, seconds int) {
	select {
	case <-ctx.Done():
		fmt.Println("in sleep func ctx error", ctx.Err())
		fmt.Println("in sleep func ctx cause", context.Cause(ctx))
	default:
		for i := 0; i < 100; i++ {
			time.Sleep(time.Second * time.Duration(seconds))
			fmt.Println("time after case", i)
		}
	}
}
