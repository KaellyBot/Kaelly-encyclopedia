package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sets"
	"github.com/rs/zerolog/log"
)

func MapSetList(dodugoSets []dodugo.SetListEntry) *amqp.EncyclopediaListAnswer {
	sets := make([]*amqp.EncyclopediaListAnswer_Item, 0)

	for _, set := range dodugoSets {
		sets = append(sets, &amqp.EncyclopediaListAnswer_Item{
			Id:   fmt.Sprintf("%v", set.GetAnkamaId()),
			Name: set.GetName(),
		})
	}

	return &amqp.EncyclopediaListAnswer{
		Items: sets,
	}
}

func MapSet(set *dodugo.EquipmentSet, items map[int32]*dodugo.Weapon,
	equipmentService equipments.Service, setService sets.Service,
) *amqp.EncyclopediaItemAnswer {
	equipments := make([]*amqp.EncyclopediaItemAnswer_Set_Equipment, 0)
	for _, itemID := range set.GetEquipmentIds() {
		formattedItemIDString := fmt.Sprintf("%v", itemID)
		item, found := items[itemID]
		if !found {
			log.Warn().
				Str(constants.LogAnkamaID, formattedItemIDString).
				Msgf("Cannot build entire set (missing item), continuing with degraded mode")
			missingItemID := itemID
			item = &dodugo.Weapon{
				AnkamaId: &missingItemID,
				Name:     &formattedItemIDString,
			}
		}

		equipments = append(equipments, &amqp.EncyclopediaItemAnswer_Set_Equipment{
			Id:    formattedItemIDString,
			Name:  item.GetName(),
			Level: int64(item.GetLevel()),
			Type:  mapEquipmentType(item.GetType(), equipmentService).EquipmentID,
		})
	}

	bonuses := make([]*amqp.EncyclopediaItemAnswer_Set_Bonus, 0)
	for _, bonus := range set.GetEffects() {
		effects := make([]*amqp.EncyclopediaItemAnswer_Effect, 0)
		var itemNumber int32
		for _, effect := range bonus {
			itemNumber = effect.GetItemCombination()
			effects = append(effects, &amqp.EncyclopediaItemAnswer_Effect{
				Id:    fmt.Sprintf("%v", *effect.GetType().Id),
				Label: effect.GetFormatted(),
			})
		}

		bonuses = append(bonuses, &amqp.EncyclopediaItemAnswer_Set_Bonus{
			ItemNumber: int64(itemNumber),
			Effects:    effects,
		})
	}

	var icon string
	setDB, found := setService.GetSetByDofusDude(set.GetAnkamaId())
	if found {
		icon = setDB.Icon
	}

	return &amqp.EncyclopediaItemAnswer{
		Type: amqp.ItemType_SET_TYPE,
		Set: &amqp.EncyclopediaItemAnswer_Set{
			Id:         fmt.Sprintf("%v", set.GetAnkamaId()),
			Name:       set.GetName(),
			Level:      int64(set.GetHighestEquipmentLevel()),
			Icon:       icon,
			IsCosmetic: set.GetIsCosmetic(),
			Equipments: equipments,
			Bonuses:    bonuses,
		},
		Source: constants.GetDofusDudeSource(),
	}
}
