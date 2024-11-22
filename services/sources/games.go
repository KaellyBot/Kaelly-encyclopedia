package sources

import (
	"context"
	"fmt"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/rs/zerolog/log"
)

func (service *Impl) ListenGameEvent(handler GameEventHandler) {
	service.eventHandlers = append(service.eventHandlers, handler)
}

func (service *Impl) checkGameVersion() {
	ctx := context.Background()
	game := amqp.Game_DOFUS_GAME
	log.Info().Msgf("Checking %v version", game)

	gameVersion, errGetDB := service.gameRepo.GetGameVersion(game)
	if errGetDB != nil {
		log.Error().Err(errGetDB).Msgf("Cannot retrieve %v version from DB, trying later...", game)
		return
	}

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

	currentVersion := gameVersion.Version
	latestGameVersion := resp.GetVersion()
	if currentVersion == latestGameVersion {
		log.Info().Msgf("No change in %v version, trying later...", game)
		return
	}

	gameVersion.Version = latestGameVersion

	if errSaveDB := service.gameRepo.Save(gameVersion); errSaveDB != nil {
		log.Error().Err(errSaveDB).Msgf("Cannot save %v version into DB, continuing...", game)
	}

	log.Info().Msgf("%v version goes from '%v' to '%v'", game, currentVersion, latestGameVersion)
	for _, handler := range service.eventHandlers {
		go emitGameEvent(handler, latestGameVersion)
	}
}

func emitGameEvent(handler GameEventHandler, gameVersion string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error().Err(fmt.Errorf("%v", err)).
				Msgf("Crash occurred while emitting game event, continuing...")
		}
	}()

	handler(gameVersion)
}
