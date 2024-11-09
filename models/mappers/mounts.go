package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
)

func MapMount(item *dodugo.Mount) *amqp.EncyclopediaItemAnswer {
	effects := make([]*amqp.EncyclopediaItemAnswer_Effect, 0)
	for _, effect := range item.GetEffects() {
		effects = append(effects, &amqp.EncyclopediaItemAnswer_Effect{
			Id:    fmt.Sprintf("%v", *effect.GetType().Id),
			Label: effect.GetFormatted(),
		})
	}

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
				ItemType:       amqp.ItemType_MOUNT_TYPE,
				EquipmentType:  amqp.EquipmentType_NONE, // TODO check with Survival, not available now
				EquipmentLabel: item.GetFamilyName(),
			},
			Icon:    *icon,
			Effects: effects,
		},
		Source: constants.GetDofusDudeSource(),
	}
}
