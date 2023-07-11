package encyclopedias

import (
	"context"
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) getSetByID(ctx context.Context, id int32, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	set, err := service.sourceService.GetSetByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	items := service.getSetEquipments(ctx, set.EquipmentIds, correlationID, lg)
	return mappers.MapSet(set, items, service.equipmentService), nil
}

func (service *Impl) getSetByQuery(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	set, err := service.sourceService.GetSetByQuery(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	items := service.getSetEquipments(ctx, set.EquipmentIds, correlationID, lg)
	return mappers.MapSet(set, items, service.equipmentService), nil
}

func (service *Impl) getSetEquipments(ctx context.Context, equipmentIds []int32, correlationID,
	lg string) map[int32]*dodugo.Weapon {
	items := make(map[int32]*dodugo.Weapon)
	for _, itemID := range equipmentIds {
		item, errItem := service.sourceService.GetEquipmentByID(ctx, itemID, lg)
		if errItem != nil {
			log.Error().Err(errItem).
				Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogAnkamaID, fmt.Sprintf("%v", itemID)).
				Msgf("Error while retrieving item with DofusDude, continuing without it")
		} else {
			items[itemID] = item
		}
	}

	return items
}
