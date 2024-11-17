package sets

import (
	"context"
	"fmt"

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

	service.sourceService.ListenGameEvent(service.buildMissingSets)
	// TODO to remove
	// service.buildMissingSets("")
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

func (service *Impl) buildMissingSets(_ string) {
	log.Info().Msgf("Building missing set icons...")
	ctx := context.Background()

	sets, errGet := service.sourceService.GetSets(ctx)
	if errGet != nil {
		log.Error().Err(errGet).Msgf("Cannot retrieve sets from DofusDude, trying later...")
		return
	}

	missingSets := make([]dodugo.SetListEntry, 0)
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
	var errorCount int
	for _, set := range missingSets {
		errBuild := service.buildMissingSet(ctx, set)
		if errBuild != nil {
			log.Warn().Err(errBuild).
				Str(constants.LogAnkamaID, fmt.Sprintf("%v", set.GetAnkamaId())).
				Msgf("Error while building set icon, continuing without this set")
			errorCount++
		}
	}

	log.Info().
		Int(constants.LogEntityCount, len(missingSets)-errorCount).
		Msg("Set icons built!")

	if errLoad := service.loadSetFromDB(); errLoad != nil {
		log.Warn().Err(errLoad).
			Msg("Could not reload set from DB, please restart to take them in account")
	}

	service.newsService.PublishSetNews(len(missingSets), len(missingSets)-errorCount)
}

func (service *Impl) buildMissingSet(ctx context.Context, set dodugo.SetListEntry,
) error {
	// Retrieve item icons
	items, errExtract := service.extractItemIcons(ctx, set)
	if errExtract != nil {
		return errExtract
	}

	// Generate set image
	image, errImg := service.buildSetImage(ctx, items)
	if errImg != nil {
		return errImg
	}

	// Write image on dedicated volume
	errUpload := writeOnDisk(set.GetAnkamaId(), image)
	if errUpload != nil {
		return errUpload
	}

	// Store cdn link into database
	errSave := service.repository.Save(entities.Set{
		DofusDudeID: set.GetAnkamaId(),
		Icon:        fmt.Sprintf(setBaseURL, set.AnkamaId),
	})
	if errSave != nil {
		return errSave
	}

	return nil
}

func (service *Impl) extractItemIcons(ctx context.Context, set dodugo.SetListEntry,
) ([]*dodugo.Weapon, error) {
	items := make([]*dodugo.Weapon, 0)
	for _, itemID := range set.GetEquipmentIds() {
		if set.GetIsCosmetic() {
			cosmetic, errCosmetic := service.sourceService.
				GetCosmeticByID(ctx, int64(itemID), constants.DofusDudeDefaultLanguage)
			if errCosmetic != nil {
				return nil, errCosmetic
			}

			items = append(items, cosmetic)
		} else {
			equipment, errItem := service.sourceService.
				GetEquipmentByID(ctx, int64(itemID), constants.DofusDudeDefaultLanguage)
			if errItem != nil {
				return nil, errItem
			}

			items = append(items, equipment)
		}
	}

	return items, nil
}
