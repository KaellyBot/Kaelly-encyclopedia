package almanaxes

import (
	"sort"
	"time"

	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/almanaxes"
	"github.com/rs/zerolog/log"
)

func New(repository repository.Repository) (*Impl, error) {
	service := Impl{
		almanaxes:  make(map[string][]entities.Almanax),
		repository: repository,
	}

	errDB := service.loadAlmanaxEffectsFromDB()
	if errDB != nil {
		return nil, errDB
	}

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
