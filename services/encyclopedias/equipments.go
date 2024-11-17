package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
)

func (service *Impl) getEquipmentByID(ctx context.Context, id int64, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	equipment, err := service.sourceService.GetEquipmentByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	ingredients := service.getIngredients(ctx, equipment.GetRecipe(), correlationID, lg)
	return mappers.MapEquipment(equipment, ingredients, service.equipmentService), nil
}

func (service *Impl) getEquipmentByQuery(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	equipment, err := service.sourceService.GetEquipmentByQuery(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	ingredients := service.getIngredients(ctx, equipment.GetRecipe(), correlationID, lg)
	return mappers.MapEquipment(equipment, ingredients, service.equipmentService), nil
}
