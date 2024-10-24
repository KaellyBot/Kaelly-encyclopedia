package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func MapEquipment(item *dodugo.Weapon, ingredientItems map[int32]*constants.Ingredient,
) *amqp.EncyclopediaItemAnswer {
	var set *amqp.EncyclopediaItemAnswer_Equipment_Set
	if item.HasParentSet() {
		parentSet := item.GetParentSet()
		set = &amqp.EncyclopediaItemAnswer_Equipment_Set{
			Id:   fmt.Sprintf("%v", parentSet.GetId()),
			Name: parentSet.GetName(),
		}
	}

	effects := make([]*amqp.EncyclopediaItemAnswer_Effect, 0)
	for _, effect := range item.GetEffects() {
		effects = append(effects, &amqp.EncyclopediaItemAnswer_Effect{
			Id:       fmt.Sprintf("%v", *effect.GetType().Id),
			Label:    effect.GetFormatted(),
			IsActive: *effect.GetType().IsActive,
		})
	}

	var recipe *amqp.EncyclopediaItemAnswer_Recipe
	if len(item.GetRecipe()) > 0 {
		ingredients := make([]*amqp.EncyclopediaItemAnswer_Recipe_Ingredient, 0)
		for _, recipeEntry := range item.GetRecipe() {
			formattedItemIDString := fmt.Sprintf("%v", recipeEntry.GetItemAnkamaId())
			ingredient, found := ingredientItems[recipeEntry.GetItemAnkamaId()]
			if !found {
				log.Warn().
					Str(constants.LogAnkamaID, formattedItemIDString).
					Msgf("Cannot build entire recipe (missing ingredient), continuing with degraded mode")
				ingredient = &constants.Ingredient{
					Name: formattedItemIDString,
					Type: amqp.ItemType_ANY_ITEM_TYPE,
				}
			}

			ingredients = append(ingredients, &amqp.EncyclopediaItemAnswer_Recipe_Ingredient{
				Id:       fmt.Sprintf("%v", recipeEntry.GetItemAnkamaId()),
				Name:     ingredient.Name,
				Quantity: int64(recipeEntry.GetQuantity()),
				Type:     ingredient.Type,
			})
		}

		recipe = &amqp.EncyclopediaItemAnswer_Recipe{
			Ingredients: ingredients,
		}
	}

	icon := item.GetImageUrls().Icon
	if item.GetImageUrls().Hq.IsSet() {
		icon = item.GetImageUrls().Hq.Get()
	}

	// TODO condition

	return &amqp.EncyclopediaItemAnswer{
		Type: amqp.ItemType_EQUIPMENT_TYPE,
		Equipment: &amqp.EncyclopediaItemAnswer_Equipment{
			Id:          fmt.Sprintf("%v", item.GetAnkamaId()),
			Name:        item.GetName(),
			Description: item.GetDescription(),
			LabelType:   *item.GetType().Name,
			Icon:        *icon,
			Level:       int64(item.GetLevel()),
			Pods:        int64(item.GetPods()),
			Set:         set,
			Effects:     effects,
			Recipe:      recipe,
		},
		Source: constants.GetDofusDudeSource(),
	}
}
