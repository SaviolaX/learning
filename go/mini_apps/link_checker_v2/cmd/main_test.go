package main

import "testing"

func TestGetInput(t *testing.T) {

	t.Run("check correct url input", func(t *testing.T) {

		url := "google.com"

		inpArgs := []string{url}

		got, _ := GetInput(inpArgs...)
		want := "google.com"

		if got != want {
			t.Errorf("want %s got %s", want, got)
		}
	})

	t.Run("check trims the spaces", func(t *testing.T) {

		url := " google.com "

		inpArgs := []string{url}

		got, _ := GetInput(inpArgs...)
		want := "google.com"

		if got != want {
			t.Errorf("want %s got %s", want, got)
		}
	})

	t.Run("check missing input error", func(t *testing.T) {

		inpArgs := []string{}

		_, err := GetInput(inpArgs...)
		want := "missing input"

		if err.Error() != want {
			t.Error("Input check failed")
		}
	})

	t.Run("check missing input error after trim", func(t *testing.T) {

		inpArgs := []string{" "}

		_, err := GetInput(inpArgs...)
		want := "input cannot be empty"

		if err.Error() != want {
			t.Error("Input check failed. Empty string was passed")
		}
	})
}

func TestFormatURL(t *testing.T) {

	t.Run("check add https:// got input url", func(t *testing.T) {

		url := "google.com"

		got, _ := FormatURL(url)
		want := "https://google.com"

		if got != want {
			t.Errorf("want %s got %s", want, got)
		}
	})

	t.Run("check error if url is empty", func(t *testing.T) {

		url := ""

		_, err := FormatURL(url)
		want := "incorrect url"

		if err.Error() != want {
			t.Error("Empty url string")
		}
	})

	t.Run("check if input url has http://", func(t *testing.T) {

		url := "https://google.com"

		got, _ := FormatURL(url)
		want := "https://google.com"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

}
