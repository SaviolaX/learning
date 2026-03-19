package url

import "context"

type Repository interface {
	Save(ctx context.Context, pair UrlPair) error
	GetByCode(ctx context.Context, code string) (UrlPair, error)
}
