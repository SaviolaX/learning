package application

import (
	"context"
	"fmt"
	dom "urlShortenerV2/internal/domain/url"
	"urlShortenerV2/internal/infrastructure/databases/json"
)

type Service struct {
	repo    dom.Repository
	hasher  func(string, int) (string, error)
	codeLen int
}

func NewService(repo dom.Repository, hasher func(string, int) (string, error), codeLen int) *Service {
	return &Service{repo: repo, hasher: hasher, codeLen: codeLen}
}

func (s *Service) Shorten(ctx context.Context, origUrl string) (string, error) {
	validatedUrl, err := dom.NewUrl(origUrl)
	if err != nil {
		return "", fmt.Errorf("validate url: %w", err)
	}

	hashCode, err := s.hasher(validatedUrl.Value, s.codeLen)
	if err != nil {
		return "", fmt.Errorf("generate hash url: %w", err)
	}

	checkedUrl, err := s.repo.GetByCode(ctx, hashCode)
	if err == nil {
		if checkedUrl.OriginUrl.Value == validatedUrl.Value {
			return hashCode, nil
		}
		return "", fmt.Errorf("hashCode collision for %s", origUrl)
	}

	pair := dom.UrlPair{Code: hashCode, OriginUrl: validatedUrl}
	if err = s.repo.Save(ctx, pair); err != nil {
		return "", fmt.Errorf("save pair: %w", err)
	}
	return hashCode, nil
}

func (s *Service) GetUrl(ctx context.Context, code string) (string, error) {

	urlPair, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return "", json.ErrNotFound
	}

	return urlPair.OriginUrl.Value, nil

}
