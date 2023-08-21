package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	ctx, c_cancel := context.WithCancelCause(context.Background())
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	go func() {
		time.Sleep(1 * time.Second)
		c_cancel(errors.New("test cancel"))
	}()

	if err := slowProcess(ctx); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Slow process completed successfully", time.Now())
	}

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("sleep 1 second and process finished", time.Now())
	case <-ctx.Done():
		fmt.Println("Process ctx done", time.Now())
	}

	fmt.Println("ctx error", ctx.Err(), time.Now())
	fmt.Println("ctx cause", context.Cause(ctx), time.Now())
}

func slowProcess(ctx context.Context) error {
	// Do some slow work here
	fmt.Println("begin in process", time.Now())
	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		// Check if the context timed out
		fmt.Println("ctx error in process", ctx.Err(), time.Now())
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("process timed out")
		}
		return ctx.Err()
	}
}
