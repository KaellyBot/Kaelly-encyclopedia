package sets

import (
	"context"
	"errors"
	"fmt"

	"github.com/dofusdude/dodugo"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

var (
	errCosmeticSet = errors.New("Set is probably a cosmetic set since no equipments could be retrieved")
)

func New(repository repository.Repository, sourceService sources.Service,
	equipmentService equipments.Service) (*Impl, error) {
	service := Impl{
		sourceService:    sourceService,
		equipmentService: equipmentService,
		sets:             make(map[int32]entities.Set),
		repository:       repository,
	}

	errDB := service.loadSetFromDB()
	if errDB != nil {
		return nil, errDB
	}

	service.sourceService.ListenGameEvent(service.buildMissingSets)
	// TODO to remove
	service.buildMissingSets()
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
	log.Info().Msgf("Building missing set icons...")
	ctx := context.Background()

	sets, errGet := service.sourceService.GetSets(ctx)
	if errGet != nil {
		log.Error().Err(errGet).Msgf("Cannot retrieve sets, trying later...")
		return
	}

	missingSets := make([]dodugo.SetListEntry, 0)
	for _, set := range sets {
		if _, found := service.sets[set.GetAnkamaId()]; !found {
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

	errLoad := service.loadSetFromDB()
	log.Warn().Err(errLoad).Msg("Could not reload set icons, please restart to take them in account")
}

func (service *Impl) buildMissingSet(ctx context.Context, set dodugo.SetListEntry) error {
	// Retrieve item icons
	items := make([]*dodugo.Weapon, 0)
	for _, itemID := range set.GetEquipmentIds() {
		item, errItem := service.sourceService.
			GetEquipmentByID(ctx, itemID, constants.DofusDudeDefaultLanguage)
		if errItem != nil {
			return errItem
		}

		if item != nil {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return errCosmeticSet
	}

	// Generate set image
	buf, errImg := service.buildSetImage(ctx, items)
	if errImg != nil {
		return errImg
	}

	// Upload image through imgur API
	imageURL, errUpload := uploadImageToImgur(ctx, buf)
	if errUpload != nil {
		return errUpload
	}

	// Store imgur link into database
	errSave := service.repository.Save(entities.Set{
		DofusDudeID: set.GetAnkamaId(),
		Icon:        imageURL,
	})
	if errSave != nil {
		return errSave
	}

	return nil
}
