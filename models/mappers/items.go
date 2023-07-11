package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
)

func MapItemList(dodugoItems []dodugo.ItemsListEntryTyped) *amqp.EncyclopediaItemListAnswer {
	items := make([]*amqp.EncyclopediaItemListAnswer_Item, 0)

	for _, item := range dodugoItems {
		items = append(items, &amqp.EncyclopediaItemListAnswer_Item{
			Id:   fmt.Sprintf("%v", item.AnkamaId),
			Name: *item.Name,
		})
	}

	return &amqp.EncyclopediaItemListAnswer{
		Items: items,
	}
}

func mapItemType(itemType dodugo.ItemsListEntryTypedType,
	equipmentService equipments.Service) amqp.EquipmentType {
	equipmentType, found := equipmentService.GetTypeByDofusDude(itemType.GetId())
	if !found {
		return amqp.EquipmentType_NONE
	}
	return equipmentType.ID
}
