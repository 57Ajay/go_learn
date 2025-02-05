package concur

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	from    string
	payload string
}

type Server struct {
	msgch chan *Message
	quit  chan struct{}
}

func (s *Server) StartAndListen(wg *sync.WaitGroup) {
	defer wg.Done()
free:
	for {
		select {
		case msg, ok := <-s.msgch:
			if !ok {
				return
			}
			// time.Sleep(10 * time.Second)
			fmt.Printf("Received message from: %s payload: %s\n", msg.from, msg.payload)
		case <-s.quit:
			fmt.Println("Server shutting down")
			break free
		}
	}
}

func sendMessageToServer(msgch chan *Message, wg *sync.WaitGroup) {
	defer wg.Done()
	msg := &Message{
		from:    "Ajay Upadhyay",
		payload: "Hello, World!",
	}
	fmt.Println("Sending single message to server")
	msgch <- msg
}

func sendManyMessagesToServer(msgch chan *Message, wg *sync.WaitGroup) {
	defer wg.Done()
	messages := []*Message{
		{from: "Ajay", payload: "Hello, World!"},
		{from: "Alice", payload: "How are you?"},
		{from: "Bob", payload: "Go is awesome!"},
	}

	fmt.Println("Sending multiple messages to server")
	for _, msg := range messages {
		time.Sleep(1 * time.Second)
		msgch <- msg
	}
}

func ChanMain() {
	channel := make(chan *Message)
	quit := make(chan struct{})
	server := &Server{
		msgch: channel,
		quit:  quit,
	}

	var wgServer sync.WaitGroup
	wgServer.Add(1)

	go server.StartAndListen(&wgServer)

	var wgSenders sync.WaitGroup
	wgSenders.Add(2)

	go sendMessageToServer(channel, &wgSenders)
	go sendManyMessagesToServer(channel, &wgSenders)

	wgSenders.Wait()
	close(channel)
	wgServer.Wait()

	for msg := range channel {
		fmt.Println("Client received:", *msg)
	}

	// Send and recieve only channels
	bidirectional := make(chan string)
	go sendOnly(bidirectional)
	go recieveOnly(bidirectional)
	fmt.Println(<-bidirectional)
}

func sendOnly(ch chan<- string) {
	ch <- "Data to send"
}

func recieveOnly(ch <-chan string) {
	data := <-ch
	fmt.Println(data)
}
