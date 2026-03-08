package main

import (
	"errors"
	"flag"
	"fmt"
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
		fmt.Println(url)

	case *fileFlag:
		file, err := GetInput(flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		fmt.Println(file)
		os.Exit(1)
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
