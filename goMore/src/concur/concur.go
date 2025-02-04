package concur

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	userId int
	name   string
}

func fetchUser(userId int, respCh chan *User, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(80 * time.Millisecond)
	respCh <- &User{userId: userId, name: "Test"}
}

func fetchUserNames(respCh chan *[]string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(69 * time.Millisecond)
	users := []User{
		{userId: 1, name: "A"},
		{userId: 2, name: "B"},
	}
	userNames := make([]string, len(users))
	for i, u := range users {
		userNames[i] = u.name
	}
	respCh <- &userNames
}

func fetchUserRecommendations(respCh chan *[]string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(69 * time.Millisecond)
	recommendations := []string{"A", "B", "C", "D"}
	respCh <- &recommendations
}

func TestMain() {
	now := time.Now()

	userChan := make(chan *User, 1)
	namesChan := make(chan *[]string, 1)
	recsChan := make(chan *[]string, 1)

	var wg sync.WaitGroup

	wg.Add(3)

	go fetchUser(21, userChan, &wg)
	go fetchUserNames(namesChan, &wg)
	go fetchUserRecommendations(recsChan, &wg)

	wg.Wait()
	fmt.Println("Waiting for response...")

	user := <-userChan
	names := <-namesChan
	recommendations := <-recsChan

	close(userChan)
	close(namesChan)
	close(recsChan)

	fmt.Println("User Details:", user.userId, user.name)

	fmt.Print("Names: ")
	for _, name := range *names {
		fmt.Print(name, " ")
	}

	fmt.Print(", Recs: ")
	for _, rec := range *recommendations {
		fmt.Print(rec, " ")
	}
	fmt.Println("\nExecution Time:", time.Since(now))
}
