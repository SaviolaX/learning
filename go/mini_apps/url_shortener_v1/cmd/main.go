package main

import (
	"fmt"
	hasher "urlShortenerV1/internal/hasher"
	url "urlShortenerV1/internal/url"
)

type UrlPair struct {
	code        string
	originalUrl string
}

const maxCodeLen = 10

func main() {
	origUrl := "youtube.com"

	validatedUrl, err := url.NewUrl(origUrl)
	if err != nil {
		fmt.Println(err)
	}

	hashCode, err := hasher.Sha256(validatedUrl.Value, maxCodeLen)
	if err != nil {
		fmt.Println(err)
	}

	pair := UrlPair{
		code:        hashCode,
		originalUrl: origUrl,
	}

	fmt.Println("validated:", pair)

}
