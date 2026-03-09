package checker

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"time"
)

type UrlReport struct {
	Url    string
	Ok     bool
	Status int
}

func CheckUrl(url string, results chan UrlReport) {
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		results <- UrlReport{Url: url, Ok: false, Status: 0}
		return
	}

	defer resp.Body.Close()

	isOk := resp.StatusCode < 300 && resp.StatusCode >= 200

	results <- UrlReport{Url: url, Ok: isOk, Status: resp.StatusCode}
}

func ReadFile(filepath string) ([]string, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return []string{}, errors.New("could not open file")
	}

	defer file.Close()

	var lstUrls []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lstUrls = append(lstUrls, line)
	}

	if err := scanner.Err(); err != nil {
		return []string{}, errors.New("scanning file error")
	}

	return lstUrls, nil

}
