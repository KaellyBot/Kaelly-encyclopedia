package sources

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/cache/v9"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func (service *Impl) getElementFromCache(ctx context.Context, key string, value any) bool {
	err := service.storeService.Get(ctx, key, value)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			log.Info().
				Str(constants.LogKey, key).
				Msgf("Cannot find element in cache, calling the API...")
		} else {
			log.Error().Err(err).
				Str(constants.LogKey, key).
				Msgf("Error while requesting element in cache, calling the API instead...")
		}
	}

	return err == nil
}

func (service *Impl) putElementToCache(ctx context.Context, key string, value any) {
	err := service.storeService.Set(ctx, key, value)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogKey, key).
			Msgf("Error while putting elment in cache, no issue to retrieve it anyway...")
	}
}

func buildListKey(objType objectType, query, language, source string) string {
	return fmt.Sprintf("%v/%v?query=%v&lg=%v", source, objType, query, language)
}

func buildItemKey(objType objectType, query, language, source string) string {
	return fmt.Sprintf("%v/%v/%v?lg=%v", source, objType, query, language)
}

func buildGameKey(source string, game amqp.Game) string {
	return fmt.Sprintf("game/%v?source=%v", game, source)
}
