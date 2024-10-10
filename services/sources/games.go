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

func (service *Impl) ListenGameEvent(handler GameEventHandler) {
	service.eventHandlers = append(service.eventHandlers, handler)
}

func (service *Impl) checkGameVersion() {
	ctx := context.Background()
	game := amqp.Game_DOFUS_GAME
	log.Info().Msgf("Checking %v version", game)
	resp, r, err := service.dofusDudeClient.MetaAPI.GetMetaVersion(ctx).Execute()
	if err != nil && r == nil {
		log.Error().Err(err).Msgf("Cannot retrieve %v version from source, trying later...", game)
		return
	}
	defer r.Body.Close()
	if err != nil {
		log.Error().Err(err).Msgf("Cannot retrieve %v version from source, trying later...", game)
		return
	}

	var gameVersion string
	key := buildGameKey(constants.GetEncyclopediasSource().Name, amqp.Game_DOFUS_GAME)
	errGet := service.storeService.Get(ctx, key, &gameVersion)
	if errGet != nil && !errors.Is(errGet, cache.ErrCacheMiss) {
		log.Error().Err(errGet).Msgf("Cannot retrieve %v version from cache, trying later...", game)
		return
	}

	latestGameVersion := resp.GetVersion()
	if gameVersion == latestGameVersion {
		log.Info().Msgf("No change in %v version, trying later...", game)
		return
	}

	if errSet := service.storeService.Set(ctx, key, latestGameVersion); errSet != nil {
		log.Error().Err(errSet).Msgf("Cannot store %v version into cache, continuing...", game)
	}

	log.Info().Msgf("%v version goes from '%v' to '%v'", game, gameVersion, latestGameVersion)
	for _, handler := range service.eventHandlers {
		emitGameEvent(handler)
	}
}

func emitGameEvent(handler GameEventHandler) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error().Err(fmt.Errorf("%v", err)).
				Msgf("Crash occurred while emitting game event, continuing...")
		}
	}()

	handler()
}
