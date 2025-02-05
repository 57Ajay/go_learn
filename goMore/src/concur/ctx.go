package concur

import (
	"context"
	"fmt"
	"time"
)

func workerWithCancel(ctx context.Context, workerID int) {
	fmt.Printf("Worker %d: Starting\n", workerID)
	defer fmt.Printf("Worker %d: Exiting\n", workerID)

	for {
		select {
		case <-ctx.Done(): // Check if context is cancelled
			fmt.Printf("Worker %d: Context cancelled, stopping\n", workerID)
			return // Exit goroutine gracefully
		default:
			fmt.Printf("Worker %d: Doing work...\n", workerID)
			time.Sleep(time.Second) // Simulate work
		}
	}
}

func doWorkWithCancel() {
	rootCtx := context.Background()            // Create a root context
	ctx, cancel := context.WithCancel(rootCtx) // Derive a cancellable context

	// Launch some worker goroutines, passing the cancellable context
	for i := 1; i <= 3; i++ {
		go workerWithCancel(ctx, i)
	}

	time.Sleep(3 * time.Second) // Let workers run for a bit
	fmt.Println("CTX: Cancelling context...")
	cancel()                    // Signal cancellation to all goroutines derived from 'ctx'
	time.Sleep(2 * time.Second) // Wait for workers to stop (graceful shutdown)

	fmt.Println("CTX: Exiting")
}

// Deriving Contexts with Deadlines and Timeouts

//  Timeout with context.WithTimeout

func workerWithTimeout(ctx context.Context, workerID int, workDuration time.Duration) {
	fmt.Printf("Worker %d: Starting, will work for %v\n", workerID, workDuration)
	defer fmt.Printf("Worker %d: Exiting\n", workerID)

	select {
	case <-time.After(workDuration): // Simulate work taking 'workDuration'
		fmt.Printf("Worker %d: Finished work successfully after %v\n", workerID, workDuration)
	case <-ctx.Done(): // Check for timeout or external cancellation
		fmt.Printf("Worker %d: Context cancelled/timeout, stopping early after %v\n", workerID, time.Since(time.Now().Add(-workDuration))) // Time spent before timeout
		return
	}
}

func doWorkWithTimeout() {
	rootCtx := context.Background()
	timeoutDuration := 2 * time.Second
	ctx, cancel := context.WithTimeout(rootCtx, timeoutDuration) // Context with timeout of 2 seconds
	defer cancel()                                               // Important: Cancel context if main function exits early to release resources

	workerWorkDuration := 3 * time.Second // Worker's intended work duration is longer than timeout

	go workerWithTimeout(ctx, 1, workerWorkDuration)

	time.Sleep(4 * time.Second) // Wait longer than timeout to see the effect

	if err := ctx.Err(); err != nil {
		fmt.Printf("Context error: %v\n", err) // Check context error after timeout
	}

	fmt.Println("TIMED_CTX: Exiting")
}

func CtxMain() {
	doWorkWithCancel()
	fmt.Println("--------------------")
	doWorkWithTimeout()
}
