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

func MoreChanMain() {
	chanBasics()
}
