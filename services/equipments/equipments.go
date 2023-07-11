package equipments

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/equipments"
	"github.com/rs/zerolog/log"
)

func New(repository repository.Repository) (*Impl, error) {
	equipmentTypes, err := repository.GetEquipmentTypes()
	if err != nil {
		return nil, err
	}

	log.Info().
		Int(constants.LogEntityCount, len(equipmentTypes)).
		Msgf("Equipment types loaded")

	dofusDudeTypes := make(map[int32]entities.EquipmentType)
	for _, equipmentType := range equipmentTypes {
		dofusDudeTypes[equipmentType.DofusDudeID] = equipmentType
	}

	return &Impl{
		dofusDudeTypes: dofusDudeTypes,
		repository:     repository,
	}, nil
}

func (service *Impl) GetTypeByDofusDude(id int32) (entities.EquipmentType, bool) {
	item, found := service.dofusDudeTypes[id]
	return item, found
}
