package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

const workers = 3

type Result struct {
	URL    string
	status int
	err    error
}

func urlToChan(urls []string) chan string {
	c := make(chan string)
	go func() {
		for _, url := range urls {
			c <- url
		}
		close(c)
	}()
	return c
}

func worker(ctx context.Context, urlsChan chan string, results chan Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case url, ok := <-urlsChan:
			if !ok {
				return
			}
			resp, err := http.Get(url)
			if err != nil {
				results <- Result{
					URL:    url,
					status: 0,
					err:    err,
				}
				continue
			}
			results <- Result{
				URL:    url,
				status: resp.StatusCode,
			}
		case <-ctx.Done():
			return
		}
	}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
		"https://anthropic.com",
		"https://youtube.com",
		"https://reddit.com",
		"https://sdkfjs.com",
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var wg sync.WaitGroup

	results := make(chan Result)
	urlsChannel := urlToChan(urls)

	for range workers {
		wg.Add(1)
		go func() {
			worker(ctx, urlsChannel, results, &wg)
		}()
	}

	go func() {
		<-quit
		cancel()
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("URL: %s\n", r.URL)
		fmt.Printf("status: %d\n", r.status)
		if r.err != nil {
			fmt.Printf("error: %v\n", r.err)
		}
	}

}
