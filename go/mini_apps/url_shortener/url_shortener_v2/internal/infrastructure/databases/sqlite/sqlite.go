package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
	dom "urlShortenerV2/internal/domain/url"
)

var ErrNotFound = errors.New("url not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	repo := &Repository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return repo, nil
}

func (r *Repository) Save(ctx context.Context, pair dom.UrlPair) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO urls (code, original) VALUES (?, ?)`,
		pair.Code, pair.OriginUrl.Value)
	if err != nil {
		return fmt.Errorf("save url: %w", err)
	}
	return nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (dom.UrlPair, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT code, original FROM urls WHERE code = ?`,
		code,
	)

	var rec struct {
		Code     string
		Original string
	}

	err := row.Scan(&rec.Code, &rec.Original)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dom.UrlPair{}, ErrNotFound
		}
		return dom.UrlPair{}, fmt.Errorf("get by code: %w", err)
	}

	u, err := dom.NewUrl(rec.Original)
	if err != nil {
		return dom.UrlPair{}, fmt.Errorf("corrupt record code=%s: %w", code, err)
	}

	return dom.UrlPair{Code: rec.Code, OriginUrl: u}, nil
}

func (r *Repository) migrate() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			code TEXT PRIMARY KEY,
			original TEXT NOT NULL
		)
	`)
	return err
}

var _ dom.Repository = (*Repository)(nil)
