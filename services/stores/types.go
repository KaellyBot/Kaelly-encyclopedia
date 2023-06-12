package stores

import (
	"context"

	"github.com/go-redis/cache/v8"
)

type Service interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any) error
}

type Impl struct {
	cache *cache.Cache
}
