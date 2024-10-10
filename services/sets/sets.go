package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func New(repository repository.Repository,
	sourceService sources.Service) (*Impl, error) {
	service := Impl{
		sourceService: sourceService,
		sets:          make(map[int32]entities.Set),
		repository:    repository,
	}

	errDB := service.loadSetFromDB()
	if errDB != nil {
		return nil, errDB
	}

	service.sourceService.ListenGameEvent(service.buildMissingSets)
	return &service, nil
}

func (service *Impl) GetSetByDofusDude(id int32) (entities.Set, bool) {
	item, found := service.sets[id]
	return item, found
}

func (service *Impl) loadSetFromDB() error {
	sets, err := service.repository.GetSets()
	if err != nil {
		return err
	}

	log.Info().
		Int(constants.LogEntityCount, len(sets)).
		Msgf("Sets loaded")

	for _, set := range sets {
		service.sets[set.DofusDudeID] = set
	}

	return nil
}

func (service *Impl) buildMissingSets() {
	// TODO sets, err := service.sourceService.GetSets()
}
