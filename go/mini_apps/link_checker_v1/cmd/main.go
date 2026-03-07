package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"linkCheckerV1/pkg"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- URL Checker ---")

	url, err := getInput("Enter your URL: ", reader)
	if err != nil {
		fmt.Println("Error reading input:", err)
	}

	urlStatus := handler.CheckUrl(url)

	fmt.Printf("\nResult for: %s\n", urlStatus.URL)
	fmt.Printf("Status Code: %d\n", urlStatus.Status)
	fmt.Printf("Is reachable: %t\n", urlStatus.Ok)
}

func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	cleaned := strings.TrimSpace(input)

	if !strings.HasPrefix(cleaned, "http") && cleaned != "" {
		cleaned = "https://" + cleaned
	}

	return cleaned, err
}
