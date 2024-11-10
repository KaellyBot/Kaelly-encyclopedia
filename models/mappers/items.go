package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
)

func MapItemList(dodugoItems []dodugo.GetGameSearch200ResponseInner) *amqp.EncyclopediaListAnswer {
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

func MapItem(item *amqp.EncyclopediaItemAnswer, language amqp.Language) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:                   amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_ANSWER,
		Status:                 amqp.RabbitMQMessage_SUCCESS,
		Language:               language,
		EncyclopediaItemAnswer: item,
	}
}

func mapEquipmentType(itemType dodugo.ItemsListEntryTypedType,
	equipmentService equipments.Service) entities.EquipmentType {
	equipmentType, found := equipmentService.GetTypeByDofusDude(itemType.GetId())
	if !found {
		return entities.EquipmentType{
			EquipmentID: amqp.EquipmentType_NONE,
			ItemID:      amqp.ItemType_ANY_ITEM_TYPE,
			DofusDudeID: itemType.GetId(),
		}
	}
	return equipmentType
}
