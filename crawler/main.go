package main

import (
	"context"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
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

func main() {
	visited := &SafeMap{urls: make(map[string]bool)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	urls := []string{
		"https://apple.com/in",
		"https://google.com",
		"https://amazon.com",
		"https://example.com",
	}

	numWorkers := 5
	urlChan := make(chan string, len(urls))
	results := make(chan string, len(urls))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, client, visited, urlChan, results)
		}()
	}

	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
