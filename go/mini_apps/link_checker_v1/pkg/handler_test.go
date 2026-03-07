package handler_test

import (
	handler "linkCheckerV1/pkg"
	"testing"
)

func TestCheckURL(t *testing.T) {
	urlOk := "https://google.com"
	urlFail := "ggl.com"

	t.Run("check working url", func(t *testing.T) {

		want := handler.URLinfo{
			URL:    urlOk,
			Ok:     true,
			Status: 200,
		}

		got := handler.CheckUrl(urlOk)

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("check not working url", func(t *testing.T) {

		want := handler.URLinfo{
			URL:    urlFail,
			Ok:     false,
			Status: 0,
		}

		got := handler.CheckUrl(urlFail)

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})
}
