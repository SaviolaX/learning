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
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func prettiPrint(report checker.UrlReport) {
	if report.Ok == true {
		fmt.Printf("|✅|%d| -> %s\n", report.Status, report.Url)
	} else {
		fmt.Printf("|❌|%d| -> %s\n", report.Status, report.Url)
	}
}

func run(args []string) error {

	fs := flag.NewFlagSet("linkCheckerV2", flag.ContinueOnError)

	urlFlag := fs.Bool("s", false, "check a single URL")
	fileFlag := fs.Bool("f", false, "check a list of URLs from file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	switch {
	case *urlFlag:
		url, err := GetInput(fs.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return err
		}

		url, err = FormatURL(url)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return err
		}

		results := make(chan checker.UrlReport)
		go checker.CheckUrl(url, results)

		r := <-results
		prettiPrint(r)

	case *fileFlag:
		filePath, err := GetInput(fs.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return err
		}
		lstUrls, err := checker.ReadFile(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return err
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
			prettiPrint(r)
		}
	}

	return nil
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
