package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
)

func (service *Impl) getMountByID(ctx context.Context, id int32, _,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	mount, err := service.sourceService.GetMountByID(ctx, id, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapMount(mount), nil
}

func (service *Impl) getMountByQuery(ctx context.Context, query, _,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	mount, err := service.sourceService.GetMountByQuery(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapMount(mount), nil
}
