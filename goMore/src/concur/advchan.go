package concur

import (
	"fmt"
	"math/rand"
	"time"
)

func dataSource(name string, ch chan string) {
	for {
		waitTime := time.Duration(rand.Intn(3)) * time.Second
		time.Sleep(waitTime)
		data := fmt.Sprintf("Data from %s at %s", name, time.Now().Format(time.RFC3339Nano))
		ch <- data
	}
}

// Waiting for Data from Multiple Channels

func DataSource() {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go dataSource("Source 1", channel1)
	go dataSource("Source 2", channel2)

	for i := 0; i < 5; i++ {
		select {
		case msg1 := <-channel1:
			fmt.Println("Received from channel 1:", msg1)
		case msg2 := <-channel2:
			fmt.Println("Received from channel 2:", msg2)
		}
	}

	fmt.Println("Exiting main")
}

//  Non-blocking Channel Operations with default

func NoBlock() {
	messageChannel := make(chan string)

	// No sender goroutine started in this example, so messageChannel will be empty

	select {
	case msg := <-messageChannel:
		fmt.Println("Received message:", msg)
	default:
		fmt.Println("No message received on messageChannel immediately (non-blocking)")
	}

	fmt.Println("Continuing in main goroutine")
}

// Timeout using select and time.After

// time.After(duration) is a handy function
// that returns a receive-only channel (<-chan Time)
// which will receive the current time after the specified
// duration has elapsed. We can use this with select to implement timeouts.

func worker(resultChan chan string) {
	time.Sleep(2 * time.Second)
	resultChan <- "Work completed!"
}

func callWorker() {
	resultChannel := make(chan string)
	go worker(resultChannel)

	timeout := 1 * time.Second

	select {
	case result := <-resultChannel:
		fmt.Println("Received result from worker:", result)
	case <-time.After(timeout):
		fmt.Println("Timeout! Worker took too long.")
	}

	fmt.Println("Exiting callWorker")
}

// Let's learn about ways to close Channels

func messageProducer(messageChannel chan string) {
	messages := []string{"Message 1", "Message 2", "Message 3", "Message 4", "Message 5"}

	for _, msg := range messages {
		messageChannel <- msg
		fmt.Println("Producer sent:", msg)
		time.Sleep(500 * time.Millisecond)
	}
	close(messageChannel) // Signal to receiver that no more messages will be sent
	fmt.Println("Producer finished sending and closed the channel")
}

func messageConsumer(messageChannel <-chan string) {
	fmt.Println("Consumer started, waiting for messages...")
	for msg := range messageChannel {
		fmt.Println("Consumer received:", msg)
	}
	fmt.Println("Consumer finished processing all messages, channel closed.")
}

func closeMain() {
	messageChan := make(chan string)

	go messageProducer(messageChan)
	go messageConsumer(messageChan)

	time.Sleep(5 * time.Second) // Long enough for all task to finish
	fmt.Println("close goroutine exiting")
}

//  Error Handling in Goroutines - panic and recover

func riskyFunction() {
	panic("Something got wrong!")
	// fmt.Println("No panic!")
}

func safeGoroutine() {
	defer func() { // Defer an anonymous function
		if r := recover(); r != nil { // Call recover() inside defer
			fmt.Println("Recovered from panic:", r) // Handle the panic
			// Optionally, we can do more error logging, cleanup, etc.
		}
	}() // immediately INVOCATION of the anonymous function

	fmt.Println("Starting safeGoroutine")
	riskyFunction()                                                      // This might panic
	fmt.Println("This line will NOT be reached if riskyFunction panics") // Won't be executed if panic occurs
	fmt.Println("Ending safeGoroutine")                                  // Won't be executed if panic occurs
}

func AdvChanMain() {
	DataSource()
	NoBlock()
	callWorker()
	closeMain()
	go safeGoroutine()
	select {} // Block forever
}
