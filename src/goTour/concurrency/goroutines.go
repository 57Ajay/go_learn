package concurrency

import (
	"fmt"
	"golang.org/x/tour/tree"
	"sync"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func worker(ch chan string) {
	ch <- "work done"
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

// Let's get grasp of sewlect indepth

func selectFibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func multipleSelect() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "Hello from c1"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "Hello from c2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println(msg1)
		case msg2 := <-c2:
			fmt.Println(msg2)
		}
	}
}

func selectWithDefault() {
	ch := make(chan int)
	select {
	case i := <-ch:
		fmt.Println("Recieved: ", i)
	default:
		fmt.Println("No message received")
	}
}

func withTimeAfter() {
	c := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		c <- "Hello from goroutine"
	}()

	select {
	case msg := <-c:
		fmt.Println(msg)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: No message received")
	}
}

func randomSelect() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		c1 <- "Hello from c1"
	}()
	go func() {
		c2 <- "Hello from c2"
	}()

	select {
	case msg := <-c1:
		fmt.Println(msg)
	case msg := <-c2:
		fmt.Println(msg)
	}
}

// Exercise: Equivalent Binary Trees

// walk walks the tree t sending all values
// from the tree to the channel ch.
func walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	walk(t.Left, ch)  // Recursively walk the left subtree
	ch <- t.Value     // Send the current node's value to the channel
	walk(t.Right, ch) // Recursively walk the right subtree
}

func same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go walk(t1, ch1)
	go walk(t2, ch2)

	for i := 0; i < 10; i++ {
		v1 := <-ch1
		v2 := <-ch2
		if v1 != v2 {
			return false
		}
	}
	return true
}

// sync.Mutex

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	return c.v[key]
}

func ConcurrencyMain() {
	fmt.Printf("THis is from concurrency package\n")
	go say("world")
	go say("hello")

	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c       // receive from c
	close(c)               // This is important to close the channel
	fmt.Println(x, y, x+y) // Now this is incredible and i am impressed

	ch := make(chan string)
	go worker(ch)
	z, ok := <-ch
	if ok {
		close(ch)
	}
	fmt.Println(z)
	fibch := make(chan int, 10)
	go fibonacci(cap(fibch), fibch)
	for i := range fibch {
		// fmt.Println(fibch)
		fmt.Println(i)
	}

	fmt.Println("----------------Select from multiple channels----------------")
	selectCh := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-selectCh)
		}
		quit <- 0
	}()
	selectFibonacci(selectCh, quit)
	multipleSelect()
	selectWithDefault()
	withTimeAfter()
	randomSelect()

	// Excercise solution
	excerch := make(chan int)
	go walk(tree.New(1), excerch)

	for i := 0; i < 10; i++ {
		fmt.Println(<-excerch) // Receive and print 10 values from the channel
	}

	fmt.Println("Are the trees same? ", same(tree.New(1), tree.New(1)))

	// sync.Mutex
	c1 := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c1.Inc("somekey")
	}

	time.Sleep(time.Second)
	fmt.Println(c1.Value("somekey"))
}
