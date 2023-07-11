package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
)

func (service *Impl) getConsumableByID(ctx context.Context, id int32, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	consumable, err := service.sourceService.GetConsumableByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	ingredients := service.getIngredients(ctx, consumable.GetRecipe(), correlationID, lg)
	return mappers.MapConsumable(consumable, ingredients), nil
}

func (service *Impl) getConsumableByQuery(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	consumable, err := service.sourceService.GetConsumableByQuery(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	ingredients := service.getIngredients(ctx, consumable.GetRecipe(), correlationID, lg)
	return mappers.MapConsumable(consumable, ingredients), nil
}
