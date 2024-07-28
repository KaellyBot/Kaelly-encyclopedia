package encyclopedias

import (
	"context"
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func (service *Impl) getIngredients(ctx context.Context, recipe []dodugo.RecipeEntry,
	correlationID, lg string) map[int32]*constants.Ingredient {
	ingredients := make(map[int32]*constants.Ingredient)
	for _, ingredient := range recipe {
		itemID := ingredient.GetItemAnkamaId()
		item, errItem := service.getIngredient(ctx, ingredient, correlationID, lg)
		if errItem != nil {
			log.Error().Err(errItem).
				Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogAnkamaID, fmt.Sprintf("%v", itemID)).
				Msgf("Error while retrieving item with DofusDude, continuing without it")
		} else {
			ingredients[itemID] = item
		}
	}

	return ingredients
}

func (service *Impl) getIngredient(ctx context.Context, ingredient dodugo.RecipeEntry,
	correlationID, lg string) (*constants.Ingredient, error) {
	itemType := service.sourceService.GetIngredientType(ingredient.GetItemSubtype())
	getIngredientByFunc, found := service.getIngredientByFuncs[itemType]
	if !found {
		return nil, sources.ErrNotFound
	}

	resp, err := getIngredientByFunc(ctx, ingredient.GetItemAnkamaId(), correlationID, lg)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Impl) getConsumableIngredientByID(ctx context.Context, id int32, correlationID,
	lg string) (*constants.Ingredient, error) {
	consumable, err := service.sourceService.GetConsumableByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	return &constants.Ingredient{
		ID:   fmt.Sprintf("%v", id),
		Name: consumable.GetName(),
		Type: amqp.IngredientType_CONSUMABLE,
	}, nil
}

func (service *Impl) getEquipmentIngredientByID(ctx context.Context, id int32, correlationID,
	lg string) (*constants.Ingredient, error) {
	equipment, err := service.sourceService.GetEquipmentByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	return &constants.Ingredient{
		ID:   fmt.Sprintf("%v", id),
		Name: equipment.GetName(),
		Type: amqp.IngredientType_EQUIPMENT_INGREDIENT,
	}, nil
}

func (service *Impl) getQuestItemIngredientByID(ctx context.Context, id int32, correlationID,
	lg string) (*constants.Ingredient, error) {
	consumable, err := service.sourceService.GetQuestItemByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	return &constants.Ingredient{
		ID:   fmt.Sprintf("%v", id),
		Name: consumable.GetName(),
		Type: amqp.IngredientType_QUEST_ITEM,
	}, nil
}

func (service *Impl) getResourceIngredientByID(ctx context.Context, id int32, correlationID,
	lg string) (*constants.Ingredient, error) {
	resource, err := service.sourceService.GetResourceByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	return &constants.Ingredient{
		ID:   fmt.Sprintf("%v", id),
		Name: resource.GetName(),
		Type: amqp.IngredientType_RESOURCE,
	}, nil
}
