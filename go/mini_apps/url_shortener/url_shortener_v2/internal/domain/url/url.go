package url

import (
	"errors"
	"fmt"
	neturl "net/url"
	"strings"
)

type Url struct {
	Value string
}

func NewUrl(originalUrl string) (Url, error) {

	originalUrl = strings.TrimSpace(originalUrl)

	if len(originalUrl) == 0 {
		return Url{}, errors.New("url is empty")
	}

	if !strings.HasPrefix(originalUrl, "http://") && !strings.HasPrefix(originalUrl, "https://") {
		originalUrl = "https://" + originalUrl
	}

	parsed, err := neturl.ParseRequestURI(originalUrl)
	if err != nil {
		return Url{}, fmt.Errorf("invalid originalUrl: %w", err)
	}

	if parsed.Host == "" {
		return Url{}, errors.New("missing host")
	}

	if !strings.Contains(parsed.Host, ".") {
		return Url{}, errors.New("invalid host")
	}

	normalized := parsed.String()

	return Url{Value: normalized}, nil

}
