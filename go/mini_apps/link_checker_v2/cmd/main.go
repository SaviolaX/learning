package main

import (
	"errors"
	"flag"
	"fmt"
	checker "linkCheckerV2/pkg"
	"os"
	"strings"
)

func main() {

	urlFlag := flag.Bool("s", false, "check a single URL")
	fileFlag := flag.Bool("f", false, "check a list of URLs from file")

	flag.Parse()

	switch {
	case *urlFlag:

		url, err := GetInput(flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		url, err = FormatURL(url)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		results := make(chan checker.UrlReport)
		go checker.CheckUrl(url, results)

		r := <-results
		fmt.Printf("URL: %s --> %d\n", r.Url, r.Status)

	case *fileFlag:
		filePath, err := GetInput(flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		lstUrls, err := checker.ReadFile(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		var formattedUrls []string
		for _, url := range lstUrls {
			formattedUrl, err := FormatURL(url)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}

			formattedUrls = append(formattedUrls, formattedUrl)
		}

		results := make(chan checker.UrlReport)
		for _, url := range formattedUrls {
			go checker.CheckUrl(url, results)
		}

		for range formattedUrls {
			r := <-results
			fmt.Printf("URL: %s --> %d\n", r.Url, r.Status)
		}
	}
}

func GetInput(args ...string) (string, error) {

	if len(args) == 0 {
		return "", errors.New("missing input")
	}

	input := strings.Join(args, " ")
	input = strings.TrimSpace(input)

	if input == "" {
		return "", errors.New("input cannot be empty")
	}

	return input, nil

}

func FormatURL(url string) (string, error) {

	if url == "" {
		return "", errors.New("incorrect url")
	}

	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		url = "https://" + url
	}

	return url, nil
}
