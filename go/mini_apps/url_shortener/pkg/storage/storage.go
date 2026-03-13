package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type URLPair struct {
	LongUrl  string
	ShortUrl string
}

type Repository struct {
	DbPath string
}

func (r *Repository) Store(urls []URLPair) error {

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

func (r *Repository) Load() ([]URLPair, error) {

	var result []URLPair

	data, err := os.ReadFile(r.DbPath)
	if err != nil {
		return result, fmt.Errorf("read file: %w", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("unmarshal storage: %w", err)
	}

	return result, nil
}
