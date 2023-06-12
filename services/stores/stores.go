package stores

import (
	"context"
	"fmt"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/spf13/viper"
)

func New() *Impl {
	return &Impl{
		cache: cache.New(&cache.Options{
			Redis: redis.NewClient(&redis.Options{
				Addr: viper.GetString(constants.RedisAddress),
			}),
			LocalCache: cache.NewTinyLFU(
				viper.GetInt(constants.RedisCacheSize),
				viper.GetDuration(constants.RedisCacheRetention),
			),
		}),
	}
}

func (service *Impl) Get(ctx context.Context, key string, value any) error {
	return service.cache.Get(ctx, buildKey(key), value)
}

func (service *Impl) Set(ctx context.Context, key string, value any) error {
	return service.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   buildKey(key),
		Value: value,
	})
}

func buildKey(query string) string {
	return fmt.Sprintf("%v/%v", constants.InternalName, query)
}
