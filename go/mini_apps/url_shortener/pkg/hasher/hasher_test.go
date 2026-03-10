package hasher

import (
	"fmt"
	"testing"
)

func TestHashUrl(t *testing.T) {

	t.Run("check empty URL error", func(t *testing.T) {
		emptyUrl := ""
		maxLen := 10

		hasher := Hasher{Url: emptyUrl}

		_, got := hasher.Sha256(maxLen)
		want := "no url passed into HashUrl"

		if got == nil {
			t.Fatal("did't get an error but wanted one")
		}

		if got.Error() != want {
			t.Errorf("got %q, wat %q", got, want)
		}
	})

	t.Run("check hashed output", func(t *testing.T) {
		url := "https://google.com"
		maxLen := 10

		hasher := Hasher{Url: url}

		got, _ := hasher.Sha256(maxLen)
		want := 10

		fmt.Println(got)

		if len(got) != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
