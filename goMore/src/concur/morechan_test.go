package concur

import (
	"fmt"
	"sync"
	"testing"
)

func TestAddUser(t *testing.T) {
	s := NewServer()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			s.AddUser(fmt.Sprintf("Ajay%d", i))
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < 1000; i++ {
		s.mu.Lock()
		if s.users[fmt.Sprintf("Ajay%d", i)] != fmt.Sprintf("Ajay%d", i) {
			s.mu.Unlock()
			t.Error("AddUser failed")
		}
		s.mu.Unlock()
	}
}
