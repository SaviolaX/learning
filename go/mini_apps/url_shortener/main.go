package main

import (
	"fmt"
	hasher "urlShortener/pkg/hasher"
)

const (
	hashLen    = 10
	defaultUrl = "https://trash-url.com/"
)

func main() {
	url := "https://google.com"
	// badUrl := ""

	hasher := hasher.Hasher{Url: url}

	hashedUrl, err := hasher.Sha256(hashLen)
	if err != nil {
		fmt.Println(err)
		return
	}

	shortUrl := defaultUrl + hashedUrl

	fmt.Println("Hashed URL:", shortUrl)

}
