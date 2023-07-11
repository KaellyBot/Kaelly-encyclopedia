package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
)

func (service *Impl) getResourceByID(ctx context.Context, id int32, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	resource, err := service.sourceService.GetResourceByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	ingredients := service.getIngredients(ctx, resource.GetRecipe(), correlationID, lg)
	return mappers.MapResource(resource, ingredients), nil
}

func (service *Impl) getResourceByQuery(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	resource, err := service.sourceService.GetResourceByQuery(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	ingredients := service.getIngredients(ctx, resource.GetRecipe(), correlationID, lg)
	return mappers.MapResource(resource, ingredients), nil
}
