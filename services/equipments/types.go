package equipments

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/equipments"
)

type Service interface {
	GetTypeByDofusDude(ID int32) (entities.EquipmentType, bool)
}

type Impl struct {
	dofusDudeTypes map[int32]entities.EquipmentType
	repository     repository.Repository
}
