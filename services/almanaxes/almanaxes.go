package almanaxes

import (
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

func (service *Impl) GetAlmanaxesByEffect(dofusDudeEffectID string) []entities.Almanax {
	return service.almanaxes[dofusDudeEffectID]
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
