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
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.load()

}


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

func (r *Repository) GetByCode(hashCode string) (UrlPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loadedUrls, err := r.load()
	if err != nil {
		return UrlPair{}, err
	}

	for _, p := range loadedUrls {
		if p.Code == hashCode {
			return p, nil
		}
	}

	return UrlPair{}, errors.New("url not found")
}
