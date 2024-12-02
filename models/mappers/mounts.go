package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
)

func MapMount(item *dodugo.Mount, equipmentService equipments.Service) *amqp.EncyclopediaItemAnswer {
	effects := make([]*amqp.EncyclopediaItemAnswer_Effect, 0)
	for _, effect := range item.GetEffects() {
		effects = append(effects, &amqp.EncyclopediaItemAnswer_Effect{
			Id:    fmt.Sprintf("%v", *effect.GetType().Id),
			Label: effect.GetFormatted(),
		})
	}

	equipmentType := mapFamilyType(item.GetFamily(), equipmentService)

	icon := item.GetImageUrls().Icon
	if item.GetImageUrls().Hq.IsSet() {
		icon = item.GetImageUrls().Hq.Get()
	}

	return &amqp.EncyclopediaItemAnswer{
		Type: amqp.ItemType_EQUIPMENT_TYPE,
		Equipment: &amqp.EncyclopediaItemAnswer_Equipment{
			Id:   fmt.Sprintf("%v", item.GetAnkamaId()),
			Name: item.GetName(),
			Type: &amqp.EncyclopediaItemAnswer_Equipment_Type{
				ItemType:       equipmentType.ItemID,
				EquipmentType:  equipmentType.EquipmentID,
				EquipmentLabel: item.Family.GetName(),
			},
			Icon:    *icon,
			Effects: effects,
		},
		Source: constants.GetDofusDudeSource(),
	}
}

func mapFamilyType(itemType dodugo.MountFamily,
	equipmentService equipments.Service) entities.EquipmentType {
	// Applying a negative ID since familyID is in conflict with equipments type;
	// These IDs were not supposed to be merged.
	mountType := itemType.GetAnkamaId() * -1
	equipmentType, found := equipmentService.GetTypeByDofusDude(mountType)
	if !found {
		return entities.EquipmentType{
			EquipmentID: amqp.EquipmentType_NONE,
			ItemID:      amqp.ItemType_MOUNT_TYPE,
			DofusDudeID: itemType.GetAnkamaId(),
		}
	}
	return equipmentType
}
