package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func MapSetList(dodugoSets []dodugo.SetListEntry) []*amqp.EncyclopediaSetListAnswer_Set {
	sets := make([]*amqp.EncyclopediaSetListAnswer_Set, 0)

	for _, set := range dodugoSets {
		sets = append(sets, &amqp.EncyclopediaSetListAnswer_Set{
			Id:   fmt.Sprintf("%v", set.GetAnkamaId()),
			Name: set.GetName(),
		})
	}

	return sets
}

func MapSet(set *dodugo.EquipmentSet, items map[int32]*dodugo.Weapon) *amqp.EncyclopediaSetAnswer {
	equipments := make([]*amqp.EncyclopediaSetAnswer_Equipment, 0)
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

		equipments = append(equipments, &amqp.EncyclopediaSetAnswer_Equipment{
			Id:   formattedItemIDString,
			Name: item.GetName(),
			Type: mapItemType(item.GetType()),
		})
	}

	bonuses := make([]*amqp.EncyclopediaSetAnswer_Bonus, 0)
	for i, bonus := range set.GetEffects() {
		effects := make([]*amqp.EncyclopediaSetAnswer_Effect, 0)
		for _, effect := range bonus {
			effects = append(effects, &amqp.EncyclopediaSetAnswer_Effect{
				Id:    fmt.Sprintf("%v", *effect.GetType().Id),
				Label: *effect.Formatted,
			})
		}

		bonuses = append(bonuses, &amqp.EncyclopediaSetAnswer_Bonus{
			ItemNumber: int64(i + constants.MinimumSetBonusItems),
			Effects:    effects,
		})
	}

	return &amqp.EncyclopediaSetAnswer{
		Id:         fmt.Sprintf("%v", set.GetAnkamaId()),
		Name:       set.GetName(),
		Level:      int64(set.GetLevel()), // TODO bug
		Equipments: equipments,
		Bonuses:    bonuses,
		Source: &amqp.Source{
			Name: constants.GetEncyclopediasSource().Name,
			Icon: constants.GetEncyclopediasSource().Icon,
			Url:  constants.GetEncyclopediasSource().URL,
		},
	}
}

func mapItemType(itemType dodugo.ItemsListEntryTypedType) amqp.EquipmentType {
	// TODO
	return amqp.EquipmentType_SHIELD
}
