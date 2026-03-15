package json

import (
	"encoding/json"
	"fmt"
	"os"
)

type UrlPair struct {
	Code        string
	OriginalUrl string
}

type Repository struct {
	DbPath string
}

func (r *Repository) Store(urls []UrlPair) error {

	// convert storage into bytes
	data, err := json.MarshalIndent(urls, "", " ")
	if err != nil {
		return fmt.Errorf("marshal storage: %w", err)
	}

	// write bytes into a file
	err = os.WriteFile(r.DbPath, data, 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (r *Repository) Load() ([]UrlPair, error) {

	var result []UrlPair

	data, err := os.ReadFile(r.DbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return result, fmt.Errorf("read file: %w", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("unmarshal storage: %w", err)
	}

	return result, nil
}

func (r *Repository) IsExist(hashCode string, pairs []UrlPair) (UrlPair, bool) {

	m := make(map[string]string)

	for _, p := range pairs {
		m[p.Code] = p.OriginalUrl
	}

	originalUrl, ok := m[hashCode]
	if ok {
		return UrlPair{Code: hashCode, OriginalUrl: originalUrl}, true

	}

	return UrlPair{}, false
}
