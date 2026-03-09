package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestGetInput(t *testing.T) {

	t.Run("correct input", func(t *testing.T) {

		input := "google.com\n"
		reader := bufio.NewReader(strings.NewReader(input))

		got, _ := GetInput("stdin", reader)
		want := "https://google.com"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("incorrect input", func(t *testing.T) {

		input := ""
		reader := bufio.NewReader(strings.NewReader(input))

		_, err := GetInput("stdin", reader)

		if err == nil {
			t.Errorf("expected error for invalid input")
		}
	})

	t.Run("check if 'https://' added", func(t *testing.T) {

		input := "google.com\n"
		reader := bufio.NewReader(strings.NewReader(input))

		got, _ := GetInput("stdin", reader)
		want := "https://google.com"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
