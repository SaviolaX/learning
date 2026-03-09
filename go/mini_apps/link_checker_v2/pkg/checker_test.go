package checker

import (
	"os"
	"testing"
	"time"
)

func TestCheckUrl(t *testing.T) {
	results := make(chan UrlReport, 1)

	t.Run("reachable url", func(t *testing.T) {

		url := "https://google.com"

		go CheckUrl(url, results)

		select {
		case res := <-results:
			if res.Url != url {
				t.Errorf("got url %s, want %s", res.Url, url)
			}

			if !res.Ok && res.Status == 0 {
				t.Errorf("failed to reach %s", url)
			}
		case <-time.After(6 * time.Second):
			t.Error("test timed out")
		}
	})

	t.Run("invalid url", func(t *testing.T) {
		url := "http://this.does.not.exist.anywhere"
		go CheckUrl(url, results)

		res := <-results
		if res.Ok {
			t.Error("expected Ok to be false for invalid domain")
		}
	})
}

func TestReadFile(t *testing.T) {
	content := "https://google.com\nhttps://github.com"
	tmpFile, err := os.CreateTemp("", "testurls.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	t.Run("correct file reading", func(t *testing.T) {
		got, err := ReadFile(tmpFile.Name())
		want := []string{"https://google.com", "https://github.com"}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(got) != len(want) || got[0] != want[0] {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("file not found error", func(t *testing.T) {
		_, err := ReadFile("non_existent.txt")
		if err == nil {
			t.Error("expected error for missing file, but got nil")
		}
	})
}
