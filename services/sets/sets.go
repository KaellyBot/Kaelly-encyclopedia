package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/rs/zerolog/log"
)

func New(repository repository.Repository) (*Impl, error) {
	sets, err := repository.GetSets()
	if err != nil {
		return nil, err
	}

	log.Info().
		Int(constants.LogEntityCount, len(sets)).
		Msgf("Sets loaded")

	dofusDudeSets := make(map[int32]entities.Set)
	for _, set := range sets {
		dofusDudeSets[set.DofusDudeID] = set
	}

	return &Impl{
		sets:       dofusDudeSets,
		repository: repository,
	}, nil
}

func (service *Impl) GetSetByDofusDude(id int32) (entities.Set, bool) {
	item, found := service.sets[id]
	return item, found
}
