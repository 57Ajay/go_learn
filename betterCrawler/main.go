package main

import (
	"context"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type SafeMap struct {
	mu   sync.Mutex
	urls map[string]bool
}

func (m *SafeMap) Visited(url string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.urls[url] {
		return true
	}
	m.urls[url] = true
	return false
}

func extractTitle(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = strings.TrimSpace(n.FirstChild.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if title == "" {
		title = "No title found"
	}
	return title, nil
}

func fetchTitle(ctx context.Context, url string, client *http.Client, visited *SafeMap, results chan<- string) {
	if visited.Visited(url) {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		results <- fmt.Sprintf("Error creating request for %s: %v", url, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		results <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return
	}

	title, err := extractTitle(resp)
	if err != nil {
		results <- fmt.Sprintf("Error parsing %s: %v", url, err)
		return
	}

	results <- fmt.Sprintf("%s - %s", url, title)
}

func worker(ctx context.Context, client *http.Client, visited *SafeMap, urls <-chan string, results chan<- string) {
	for url := range urls {
		fetchTitle(ctx, url, client, visited, results)
	}
}

func getURLFromUser() string {
	fmt.Println("\nWebsite Title Crawler")
	fmt.Println("=====================")
	fmt.Println("Enter a website URL to get its title")
	fmt.Println("Examples:")
	fmt.Println("- apple.com/in      (will become https://apple.com/in)")
	fmt.Println("- https://google.com")
	fmt.Println("- http://example.org")
	fmt.Print("\nYour URL: ")

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	// Add scheme if missing
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "https://" + input
	}

	// Basic validation
	if !strings.Contains(input, ".") || len(input) < 10 {
		fmt.Println("Invalid URL format. Please include at least a domain name.")
		fmt.Println("Example: apple.com/in or https://example.org")
		os.Exit(1)
	}

	return input
}

func main() {
	visited := &SafeMap{urls: make(map[string]bool)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Do not follow redirects
		},
	}

	userURL := getURLFromUser()

	urlChan := make(chan string, 1)
	results := make(chan string, 1)

	// Start worker
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker(ctx, client, visited, urlChan, results)
	}()

	urlChan <- userURL
	close(urlChan)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("\nResult:")
	if result, ok := <-results; ok {
		fmt.Println(result)
	} else {
		fmt.Println("No result received (timeout or error)")
	}
}
