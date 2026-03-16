package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

var ErrNotFound = errors.New("url not found")

type UrlPair struct {
	Code        string
	OriginalUrl string
}

type Repository struct {
	DbPath string
	mu     sync.RWMutex
}

func (r *Repository) Store(pair UrlPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pairs, err := r.load()
	if err != nil {
		return err
	}

	pairs = append(pairs, pair)
	return r.write(pairs)
}

func (r *Repository) write(urls []UrlPair) error {
	data, err := json.MarshalIndent(urls, "", " ")
	if err != nil {
		return fmt.Errorf("marshal storage: %w", err)
	}
	return os.WriteFile(r.DbPath, data, 0644)
}

func (r *Repository) Load() ([]UrlPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.load()

}

func (r *Repository) load() ([]UrlPair, error) {

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

	return UrlPair{}, ErrNotFound
}
