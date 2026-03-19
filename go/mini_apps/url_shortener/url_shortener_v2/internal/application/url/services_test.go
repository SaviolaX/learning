package application

import (
	"context"
	"path/filepath"
	"testing"
	dom "urlShortenerV2/internal/domain/url"
	jsonRepo "urlShortenerV2/internal/infrastructure/databases/json"
	"urlShortenerV2/internal/infrastructure/hasher"

	"github.com/stretchr/testify/assert"
)

func TestShorten_InvalidUrl(t *testing.T) {
	repo := newTestRepo(t)

	service := NewService(repo, hasher.Sha256, 10)

	_, err := service.Shorten(context.Background(), "not_a_url")

	assert.Error(t, err)

}

func TestShorten_UrlAlreadyExists(t *testing.T) {
	repo := newTestRepo(t)

	service := NewService(repo, hasher.Sha256, 10)
	origUrl := "https://google.com"

	r, err := service.Shorten(context.Background(), origUrl)

	assert.NoError(t, err)
	assert.Len(t, r, 10)

	res, err := service.Shorten(context.Background(), origUrl)

	assert.NoError(t, err)
	assert.Equal(t, r, res)

}

func TestShorten_Correct(t *testing.T) {
	repo := newTestRepo(t)

	service := NewService(repo, hasher.Sha256, 10)
	origUrl := "https://google.com"

	res, err := service.Shorten(context.Background(), origUrl)

	assert.NoError(t, err)
	assert.Len(t, res, 10)

}

func TestGetUrl_Found(t *testing.T) {
	repo := newTestRepo(t)

	pair := dom.UrlPair{Code: "asdf", OriginUrl: dom.Url{Value: "https://google.com"}}
	repo.Save(context.Background(), pair)

	service := NewService(repo, hasher.Sha256, 10)

	got, err := service.GetUrl(context.Background(), pair.Code)
	want := pair.OriginUrl.Value

	assert.NoError(t, err)
	assert.Equal(t, want, got)

}

func TestGetUrl_NotFound(t *testing.T) {
	repo := newTestRepo(t)

	service := NewService(repo, hasher.Sha256, 10)

	res, err := service.GetUrl(context.Background(), "asdf")

	assert.ErrorIs(t, err, jsonRepo.ErrNotFound)
	assert.Equal(t, "", res)

}

func newTestRepo(t *testing.T) *jsonRepo.Repository {
	t.Helper()
	dir := t.TempDir()
	return jsonRepo.NewRepository(filepath.Join(dir, "db.json"))
}
