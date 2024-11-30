package sets

import (
	"context"

	"github.com/dofusdude/dodugo"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/news"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func New(repository repository.Repository, newsService news.Service,
	sourceService sources.Service, equipmentService equipments.Service) (*Impl, error) {
	service := Impl{
		newsService:      newsService,
		sourceService:    sourceService,
		equipmentService: equipmentService,
		sets:             make(map[int64]entities.Set),
		repository:       repository,
	}

	errDB := service.loadSetFromDB()
	if errDB != nil {
		return nil, errDB
	}

	service.sourceService.ListenGameEvent(service.checkMissingSets)
	return &service, nil
}

func (service *Impl) GetSetByDofusDude(id int64) (entities.Set, bool) {
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
		service.sets[int64(set.DofusDudeID)] = set
	}

	return nil
}

func (service *Impl) checkMissingSets(_ string) {
	log.Info().Msgf("Checking missing set icons...")
	ctx := context.Background()

	sets, errGet := service.sourceService.GetSets(ctx)
	if errGet != nil {
		log.Error().Err(errGet).Msgf("Cannot retrieve sets from DofusDude, trying later...")
		return
	}

	missingSets := make([]dodugo.ListEquipmentSet, 0)
	for _, set := range sets {
		if _, found := service.sets[int64(set.GetAnkamaId())]; !found {
			missingSets = append(missingSets, set)
		}
	}

	if len(missingSets) == 0 {
		log.Info().Int(constants.LogEntityCount, len(missingSets)).Msgf("Set icons are all up-to-date")
		return
	}

	log.Info().Int(constants.LogEntityCount, len(missingSets)).Msgf("Set icons to build")
	service.newsService.PublishSetNews(missingSets)
}
