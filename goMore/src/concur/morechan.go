package concur

import (
	"fmt"
	"sync"
)

func chanBasics() {
	strchan := make(chan string, 4)
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		strchan <- "Ajay Upadhyay"
		wg.Done()
	}()
	go func() {
		strchan <- "Ajay Upadhyay is "
		wg.Done()
	}()

	go func() {
		strchan <- "Ajay Upadhyay is coming"
		wg.Done()
	}()
	go func() {
		strchan <- "Ajay Upadhyay is coming here"
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(strchan)
	}()

	for val := range strchan {
		fmt.Println("from MoreChanMain: ", val)
	}
	fmt.Println("Channel Closed and Range Exited")
}

type Server_chan struct {
	users map[string]string
	mu    sync.Mutex
}

func NewServer() *Server_chan {
	return &Server_chan{users: make(map[string]string)}
}

func (s *Server_chan) AddUser(user string) {
	s.mu.Lock()
	s.users[user] = user
	s.mu.Unlock()
}

func MoreChanMain() {
	chanBasics()
}
