package main

import (
	"fmt"
	"net/http"
	//	"sync"
	"time"
)

type URLinfo struct {
	URL string
	OK  bool
}

func main() {

	urls := []string{
		"https://gogle.com",  // not work
		"https://google.com", // work
		"https://youtube.com",
		"https://github.com",
	}

	client := http.Client{Timeout: 5 * time.Second}

	results := make(chan URLinfo)

	// Use channels
	for _, url := range urls {
		go checkUrl(url, client, results)

	}

	for range urls {
		r := <-results
		fmt.Println(r.URL, "->", r.OK)
	}

	// Use Sleep to waint gorutines are done
	// for _, url := range urls {
	// 	go checkUrl(url, client)
	// }
	// time.Sleep(5 * time.Second)

	// Use WaitGroup for gorutines
	//var wg sync.WaitGroup

	//for _, url := range urls {

	//	wg.Add(1)

	//	go func(url string) {
	//		defer wg.Done()
	//		checkUrl(url, client)
	//	}(url)

	//}

	//wg.Wait()
}

func checkUrl(url string, client http.Client, results chan URLinfo) {
	resp, err := client.Get(url)
	if err != nil {
		results <- URLinfo{URL: url, OK: false}
		return
	}
	results <- URLinfo{URL: url, OK: true}

	resp.Body.Close()
}
