package mappers

import (
	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
)

func mapItemType(itemType dodugo.ItemsListEntryTypedType,
	equipmentService equipments.Service) amqp.EquipmentType {
	equipmentType, found := equipmentService.GetTypeByDofusDude(itemType.GetId())
	if !found {
		return amqp.EquipmentType_NONE
	}
	return equipmentType.ID
}
