package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
)

func MapItemList(dodugoItems []dodugo.ItemsListEntryTyped) *amqp.EncyclopediaListAnswer {
	items := make([]*amqp.EncyclopediaListAnswer_Item, 0)

	for _, item := range dodugoItems {
		items = append(items, &amqp.EncyclopediaListAnswer_Item{
			Id:   fmt.Sprintf("%v", item.AnkamaId),
			Name: *item.Name,
		})
	}

	return &amqp.EncyclopediaListAnswer{
		Items: items,
	}
}

func MapQuestItem(item *dodugo.Resource, ingredientItems map[int32]*constants.Ingredient,
) *amqp.EncyclopediaItemAnswer {
	// TODO

	return &amqp.EncyclopediaItemAnswer{
		Type:      amqp.ItemType_QUEST_ITEM,
		QuestItem: &amqp.EncyclopediaItemAnswer_QuestItem{},
		Source: constants.GetDofusDudeSource(),
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
