// Concurrency patterns are reusable, well-established
// solutions to common problems that arise in concurrent
// programming. They are like design patterns for concurrency,
// providing blueprints for structuring your concurrent code
// to be efficient, reliable, and easier to understand.
// Instead of reinventing the wheel every time you face
// a concurrent problem, you can often apply a known pattern.

package concur

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID      int
	Payload int
}

func worker_(workerID int, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done() // Signal worker completion

	fmt.Printf("Worker %d: Starting\n", workerID)
	for task := range tasks { // Range over tasks channel until it's closed
		fmt.Printf("Worker %d: Received task ID %d, Payload: %d\n", workerID, task.ID, task.Payload)
		time.Sleep(time.Duration(task.Payload) * time.Millisecond) // Simulate task processing time
		fmt.Printf("Worker %d: Completed task ID %d\n", workerID, task.ID)
	}
	fmt.Printf("Worker %d: Shutting down\n", workerID)
}

func doWork() {
	numWorkers := 3                           // Number of worker goroutines in the pool
	numTasks := 10                            // Total number of tasks to submit
	tasksChannel := make(chan Task, numTasks) // Buffered channel for tasks
	var wg sync.WaitGroup                     // WaitGroup to wait for all workers to finish

	// Launch worker goroutines
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker_(i, tasksChannel, &wg)
	}

	// Submit tasks to the tasks channel
	for i := 1; i <= numTasks; i++ {
		task := Task{ID: i, Payload: i * 100} // Example task with ID and payload
		tasksChannel <- task
		fmt.Printf("Main: Submitted task ID %d\n", task.ID)
	}
	close(tasksChannel) // Signal to workers that no more tasks are coming

	fmt.Println("Main: Waiting for workers to finish...")
	wg.Wait() // Wait for all workers to complete
	fmt.Println("Main: All workers finished, exiting")
}

func PatrnMain() {
	doWork()
}
