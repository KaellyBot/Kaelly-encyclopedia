package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func MapItemList(dodugoItems []dodugo.ItemsListEntryTyped) []*amqp.EncyclopediaItemListAnswer_Item {
	items := make([]*amqp.EncyclopediaItemListAnswer_Item, 0)

	for _, item := range dodugoItems {
		items = append(items, &amqp.EncyclopediaItemListAnswer_Item{
			Id:   fmt.Sprintf("%v", item.AnkamaId),
			Name: *item.Name,
		})
	}

	return items
}

func MapItem(item *dodugo.Weapon, ingredientItems map[int32]*dodugo.Weapon,
) *amqp.EncyclopediaItemAnswer {

	// TODO the rest

	var set *amqp.EncyclopediaItemAnswer_Set
	if item.HasParentSet() {
		parentSet := item.GetParentSet()
		set = &amqp.EncyclopediaItemAnswer_Set{
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
		ingredients := make([]*amqp.EncyclopediaItemAnswer_Ingredient, 0)
		for _, recipeEntry := range item.GetRecipe() {
			formattedItemIDString := fmt.Sprintf("%v", recipeEntry.GetItemAnkamaId())
			ingredient, found := ingredientItems[recipeEntry.GetItemAnkamaId()]
			if !found {
				log.Warn().
					Str(constants.LogAnkamaID, formattedItemIDString).
					Msgf("Cannot build entire recipe (missing ingredient), continuing with degraded mode")
				missingIngredientID := recipeEntry.GetItemAnkamaId()
				ingredient = &dodugo.Weapon{
					AnkamaId: &missingIngredientID,
					Name:     &formattedItemIDString,
				}
			}
			ingredients = append(ingredients, &amqp.EncyclopediaItemAnswer_Ingredient{
				Id:       fmt.Sprintf("%v", recipeEntry.GetItemAnkamaId()),
				Name:     ingredient.GetName(),
				Quantity: int64(recipeEntry.GetQuantity()),
				// TODO bring also type to determine the URL
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

	return &amqp.EncyclopediaItemAnswer{
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
		Source: &amqp.Source{
			Name: constants.GetEncyclopediasSource().Name,
			Icon: constants.GetEncyclopediasSource().Icon,
			Url:  constants.GetEncyclopediasSource().URL,
		},
	}
}
