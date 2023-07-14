package mappers

import (
	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
)

func MapCosmetic(item *dodugo.Cosmetic, ingredientItems map[int32]*constants.Ingredient,
) *amqp.EncyclopediaItemAnswer {
	// TODO

	return &amqp.EncyclopediaItemAnswer{
		Type:     amqp.ItemType_COSMETIC,
		Cosmetic: &amqp.EncyclopediaItemAnswer_Cosmetic{},
		Source: &amqp.Source{
			Name: constants.GetEncyclopediasSource().Name,
			Icon: constants.GetEncyclopediasSource().Icon,
			Url:  constants.GetEncyclopediasSource().URL,
		},
	}
}
