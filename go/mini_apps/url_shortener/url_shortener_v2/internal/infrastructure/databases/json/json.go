package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	dom "urlShortenerV2/internal/domain/url"
)

var ErrNotFound = errors.New("url not found")

type urlRecord struct {
	Code      string
	OriginUrl string
}

type Repository struct {
	dbPath string
	mu     sync.RWMutex
}

func NewRepository(dbPath string) *Repository {
	return &Repository{
		dbPath: dbPath,
	}
}

func (r *Repository) Save(ctx context.Context, pair dom.UrlPair) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	pairs, err := r.load()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	pairs = append(pairs, toRecord(pair))
	return r.write(pairs)
}

func (r *Repository) GetByCode(ctx context.Context, hashCode string) (dom.UrlPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loadedUrls, err := r.load()
	if err != nil {
		return dom.UrlPair{}, err
	}

	for _, p := range loadedUrls {
		if p.Code == hashCode {
			return toDomain(p), nil
		}
	}

	return dom.UrlPair{}, ErrNotFound
}

func (r *Repository) write(urls []urlRecord) error {
	data, err := json.MarshalIndent(urls, "", " ")
	if err != nil {
		return fmt.Errorf("marshal storage: %w", err)
	}
	return os.WriteFile(r.dbPath, data, 0644)
}

func (r *Repository) load() ([]urlRecord, error) {

	data, err := os.ReadFile(r.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	var records []urlRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("unmarshal storage: %w", err)
	}

	return records, nil
}

func toDomain(r urlRecord) dom.UrlPair {
	return dom.UrlPair{
		Code:      r.Code,
		OriginUrl: dom.Url{Value: r.OriginUrl},
	}
}

func toRecord(pair dom.UrlPair) urlRecord {
	return urlRecord{
		Code:      pair.Code,
		OriginUrl: pair.OriginUrl.Value,
	}
}
