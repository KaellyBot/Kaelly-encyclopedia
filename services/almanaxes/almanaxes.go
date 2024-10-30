package almanaxes

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/almanaxes"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func New(repository repository.Repository,
	sourceService sources.Service) (*Impl, error) {
	service := Impl{
		almanaxes:     make(map[string][]entities.Almanax),
		sourceService: sourceService,
		repository:    repository,
	}

	errDB := service.loadAlmanaxEffectsFromDB()
	if errDB != nil {
		return nil, errDB
	}

	service.sourceService.ListenGameEvent(service.reconcileDofusDudeIDs)

	return &service, nil
}

func (service *Impl) GetDatesByAlmanaxEffect(dofusDudeEffectID string) []time.Time {
	now := time.Now().UTC()
	dates := make([]time.Time, 0)
	entities, found := service.almanaxes[dofusDudeEffectID]
	if !found {
		return dates
	}

	for _, entity := range entities {
		year := now.Year()
		if time.Month(entity.Month) < now.Month() ||
			time.Month(entity.Month) == now.Month() && entity.Day < now.Day() {
			year++
		}

		date := time.Date(year, time.Month(entity.Month), entity.Day, 0, 0, 0, 0, time.UTC)
		dates = append(dates, date)
	}

	// Sorted by ASC
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	return dates
}

func (service *Impl) loadAlmanaxEffectsFromDB() error {
	almanaxes, err := service.repository.GetAlmanaxes()
	if err != nil {
		return err
	}

	log.Info().
		Int(constants.LogEntityCount, len(almanaxes)).
		Msgf("Almanaxes loaded")

	for _, almanax := range almanaxes {
		effects := service.almanaxes[almanax.DofusDudeEffectID]
		service.almanaxes[almanax.DofusDudeEffectID] = append(effects, almanax)
	}

	return nil
}

func (service *Impl) reconcileDofusDudeIDs() {
	log.Info().Msgf("Reconciling almanax DofusDude IDs...")
	ctx := context.Background()
	year := time.Now().Year()

	almanaxEntities, errDB := service.repository.GetAlmanaxes()
	if errDB != nil {
		log.Error().Err(errDB).Msgf("Cannot retrieve almanaxes from database, trying later...")
		return
	}

	var updatedCount int
	var errorCount int
	for _, almanaxEntity := range almanaxEntities {
		updated, errRec := service.reconcileDofusDudeID(ctx, almanaxEntity, year)
		if errRec != nil {
			log.Warn().Err(errRec).
				Str(constants.LogDate, fmt.
					Sprintf("%v-%v-%v", year, almanaxEntity.Month, almanaxEntity.Day)).
				Msgf("Error while reconciliating almanax, continuing without this date")
			errorCount++
			continue
		}

		if updated {
			updatedCount++
		}
	}

	if updatedCount == 0 {
		log.Info().
			Int(constants.LogEntityCount, len(almanaxEntities)).
			Msgf("Almanax days are all up-to-date")
		return
	}

	log.Info().
		Int(constants.LogEntityCount, updatedCount).
		Msg("Almanax dates reconciliated!")

	errLoad := service.loadAlmanaxEffectsFromDB()
	log.Warn().Err(errLoad).Msg("Could not reload almanax from DB, please restart to take them in account")
}

func (service *Impl) reconcileDofusDudeID(ctx context.Context, entity entities.Almanax, year int,
) (bool, error) {
	day := time.Date(year, time.Month(entity.Month), entity.Day, 0, 0, 0, 0, time.UTC)
	dodugoAlmanax, errGet := service.sourceService.
		GetAlmanaxByDate(ctx, day, constants.DofusDudeDefaultLanguage)
	if errGet != nil {
		return false, errGet
	}

	if dodugoAlmanax == nil {
		return false, errNotFound
	}

	if dodugoAlmanax.Bonus.Type.GetId() != entity.DofusDudeEffectID {
		entity.DofusDudeEffectID = dodugoAlmanax.Bonus.Type.GetId()
		return true, service.repository.Save(entity)
	}

	return false, nil
}
