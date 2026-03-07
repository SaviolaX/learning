package handler

import (
	"net/http"
	"time"
)

type URLinfo struct {
	URL    string
	Ok     bool
	Status int
}

func CheckUrl(url string) URLinfo {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return URLinfo{URL: url, Ok: false, Status: 0}
	}

	defer resp.Body.Close()

	isOk := resp.StatusCode >= 200 && resp.StatusCode < 300

	return URLinfo{URL: url, Ok: isOk, Status: resp.StatusCode}
}
