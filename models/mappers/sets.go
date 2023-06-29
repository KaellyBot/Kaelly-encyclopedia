package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
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

func MapSet(set *dodugo.EquipmentSet) *amqp.EncyclopediaSetAnswer {
	equipments := make([]*amqp.EncyclopediaSetAnswer_Equipment, 0)
	for _, equipmentId := range set.GetEquipmentIds() {
		equipments = append(equipments, &amqp.EncyclopediaSetAnswer_Equipment{
			Id:   fmt.Sprintf("%v", equipmentId),
			Name: fmt.Sprintf("%v", equipmentId), // TODO
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
		Level:      int64(set.GetLevel()), // TODO bug, does not work
		Equipments: equipments,
		Bonuses:    bonuses,
		Source: &amqp.EncyclopediaSetAnswer_Source{
			Name: constants.GetEncyclopediasSource().Name,
			Icon: constants.GetEncyclopediasSource().Icon,
			Url:  constants.GetEncyclopediasSource().URL,
		},
	}
}
