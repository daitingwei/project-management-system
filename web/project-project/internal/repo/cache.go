package repo

import (
	"context"
	"time"
)

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HSet(ctx context.Context, key string, field string, value string)
	HKeys(ctx context.Context, key string) ([]string, error)
	Delete(background context.Context, keys []string)
	Keys(ctx context.Context, pattern string) ([]string, error)
}
