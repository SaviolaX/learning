package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Url struct {
	value string
}

type UrlPair struct {
	code        string
	originalUrl string
}

func main() {
	orgUrl := "sdfja;ls"

	validatedUrl, err := NewUrl(orgUrl)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("validated:", validatedUrl)

}

func NewUrl(originalUrl string) (Url, error) {

	if len(originalUrl) == 0 {
		return Url{}, errors.New("url is empty")
	}

	originalUrl = strings.TrimSpace(originalUrl)

	if !strings.HasPrefix(originalUrl, "http://") && !strings.HasPrefix(originalUrl, "https://") {
		originalUrl = "https://" + originalUrl
	}

	parsed, err := url.ParseRequestURI(originalUrl)
	if err != nil {
		return Url{}, fmt.Errorf("invalid originalUrl: %w", err)
	}

	if parsed.Host == "" {
		return Url{}, errors.New("missing host")
	}

	if !strings.Contains(parsed.Host, ".") {
		return Url{}, errors.New("invalid host")
	}

	normalized := parsed.String()

	return Url{value: normalized}, nil

}
