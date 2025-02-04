package concur

import (
	"fmt"
	"sync"
)

type Message struct {
	from    string
	payload string
}

type Server struct {
	msgch chan *Message
}

func (s *Server) StartAndListen(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range s.msgch {
		fmt.Printf("Received message from: %s payload: %s\n", msg.from, msg.payload)
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
		msgch <- msg
	}
}

func ChanMain() {
	channel := make(chan *Message)

	server := &Server{
		msgch: channel,
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
}
