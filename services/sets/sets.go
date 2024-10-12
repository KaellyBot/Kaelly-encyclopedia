package sets

import (
	"context"
	"fmt"

	"github.com/dofusdude/dodugo"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
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
	for _, set := range missingSets {
		service.buildMissingSet(ctx, set)
	}
}

func (service *Impl) buildMissingSet(ctx context.Context, set dodugo.SetListEntry) {
	// Retrieve item icons
	items := make([]*dodugo.Weapon, 0)
	for _, itemID := range set.GetEquipmentIds() {
		item, errItem := service.sourceService.
			GetEquipmentByID(ctx, itemID, constants.DofusDudeDefaultLanguage)
		if errItem != nil {
			log.Warn().Err(errItem).
				Str(constants.LogAnkamaID, fmt.Sprintf("%v", set.GetAnkamaId())).
				Msgf("Error while retrieving item with DofusDude, continuing without this set")
			return
		}

		items = append(items, item)
	}

	// Generate set image
	buf, errImg := service.buildSetImage(ctx, items)
	if errImg != nil {
		log.Warn().Err(errImg).
			Str(constants.LogAnkamaID, fmt.Sprintf("%v", set.GetAnkamaId())).
			Msgf("Error while generating set icon, continuing without this set")
		return
	}

	// Upload image through imgur API
	imageURL, errUpload := uploadImageToImgur(ctx, buf)
	if errUpload != nil {
		log.Warn().Err(errUpload).
			Str(constants.LogAnkamaID, fmt.Sprintf("%v", set.GetAnkamaId())).
			Msgf("Error while uploading set icon, continuing without this set")
		return
	}

	// Store imgur link into database
	errSave := service.repository.Save(entities.Set{
		DofusDudeID: set.GetAnkamaId(),
		Icon:        imageURL,
	})
	if errSave != nil {
		log.Warn().Err(errSave).
			Str(constants.LogAnkamaID, fmt.Sprintf("%v", set.GetAnkamaId())).
			Msgf("Error while saving set icon, continuing without this set")
		return
	}
}
