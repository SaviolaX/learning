package sqlite

import (
	"context"
	// "path/filepath"
	"testing"
	dom "urlShortenerV2/internal/domain/url"

	"github.com/stretchr/testify/assert"
)

func TestGetByCode_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	_, err := repo.GetByCode(context.Background(), "doesnotexist")

	assert.ErrorIs(t, err, ErrNotFound)

}

func TestGetByCode_Correct(t *testing.T) {
	repo := newTestRepo(t)
	pair := dom.UrlPair{Code: "asdf", OriginUrl: dom.Url{Value: "https://google.com"}}
	repo.Save(context.Background(), pair)

	res, err := repo.GetByCode(context.Background(), pair.Code)

	assert.NoError(t, err)
	assert.Equal(t, pair, res)

}

func TestSave_Correct(t *testing.T) {
	repo := newTestRepo(t)

	pair := dom.UrlPair{Code: "asdf", OriginUrl: dom.Url{Value: "https://google.com"}}

	err := repo.Save(context.Background(), pair)

	assert.NoError(t, err)

}

func TestSave_InvalidPath(t *testing.T) {
	_, err := NewRepository("this/is/incorrect/path.db")

	assert.Error(t, err)
}

func newTestRepo(t *testing.T) *Repository {
	t.Helper()
	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create test repository: %v", err)
	}

	return repo
}
