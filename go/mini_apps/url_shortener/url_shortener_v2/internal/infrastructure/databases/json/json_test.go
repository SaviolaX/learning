package json

import (
	"context"
	"path/filepath"
	"testing"
	dom "urlShortenerV2/internal/domain/url"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryGetByCode_Correct(t *testing.T) {
	repo := newTestRepo(t)

	pair := dom.UrlPair{
		Code:      "asdf",
		OriginUrl: dom.Url{Value: "https://google.com"},
	}

	repo.Save(context.Background(), pair)

	result, err := repo.GetByCode(context.Background(), pair.Code)

	assert.NoError(t, err)
	assert.Equal(t, pair, result)
}

func TestRepositoryGetByCode_NotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.GetByCode(context.Background(), "doesnotexits")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRepositorySave_Correct(t *testing.T) {

	repo := newTestRepo(t)

	pair := dom.UrlPair{
		Code:      "asdf",
		OriginUrl: dom.Url{Value: "https://google.com"},
	}

	err := repo.Save(context.Background(), pair)

	assert.NoError(t, err)
}

func newTestRepo(t *testing.T) *Repository {
	t.Helper()
	dir := t.TempDir()
	return &Repository{dbPath: filepath.Join(dir, "db.json")}
}

func TestRepositorySave_IncorrectPath(t *testing.T) {
	repo := Repository{dbPath: "incorrect/path/db.json"}

	pair := dom.UrlPair{
		Code:      "asdf",
		OriginUrl: dom.Url{Value: "https://google.com"},
	}

	err := repo.Save(context.Background(), pair)

	assert.Error(t, err)

}
