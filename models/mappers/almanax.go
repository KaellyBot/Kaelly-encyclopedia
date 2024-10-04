package mappers

import (
	"fmt"
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapAlmanax(dodugoAlmanax *dodugo.AlmanaxEntry) *amqp.Almanax {
	if dodugoAlmanax == nil {
		return nil
	}

	date, err := time.Parse(constants.DofusDudeAlmanaxDateFormat, *dodugoAlmanax.Date)
	if err != nil {
		log.Warn().
			Str(constants.LogDate, *dodugoAlmanax.Date).
			Msgf("Cannot cast dofusdude almanax date, continuing with time.Now...")
		date = time.Now()
	}

	icon := *dodugoAlmanax.Tribute.Item.GetImageUrls().Icon
	if dodugoAlmanax.Tribute.Item.GetImageUrls().Sd.IsSet() {
		icon = *dodugoAlmanax.Tribute.Item.GetImageUrls().Sd.Get()
	}

	return &amqp.Almanax{
		Bonus: *dodugoAlmanax.Bonus.Description,
		Tribute: &amqp.Almanax_Tribute{
			Item: &amqp.Almanax_Tribute_Item{
				Name: *dodugoAlmanax.Tribute.Item.Name,
				Icon: icon,
			},
			Quantity: *dodugoAlmanax.Tribute.Quantity,
		},
		Reward: int64(dodugoAlmanax.GetRewardKamas()),
		Date:   timestamppb.New(date.UTC()),
		Source: constants.GetDofusDudeSource(),
	}
}

func MapAlmanaxEffectList(dodugoAlmanaxEffects []dodugo.GetMetaAlmanaxBonuses200ResponseInner,
) *amqp.EncyclopediaListAnswer {
	effects := make([]*amqp.EncyclopediaListAnswer_Item, 0)

	for _, effect := range dodugoAlmanaxEffects {
		effects = append(effects, &amqp.EncyclopediaListAnswer_Item{
			Id:   fmt.Sprintf("%v", effect.GetId()),
			Name: effect.GetName(),
		})
	}

	return &amqp.EncyclopediaListAnswer{
		Items: effects,
	}
}

func MapAlmanaxResource(dodugoAlmanax []dodugo.AlmanaxEntry, dayDuration int32,
) *amqp.EncyclopediaAlmanaxResourceAnswer {
	resources := make(map[string]int32, 0)
	for _, almanax := range dodugoAlmanax {
		itemName := *almanax.Tribute.GetItem().Name
		quantity, found := resources[itemName]
		if !found {
			quantity = 0
		}

		resources[itemName] = quantity + almanax.Tribute.GetQuantity()
	}

	tributes := make([]*amqp.EncyclopediaAlmanaxResourceAnswer_Tribute, 0)
	for itemName, quantity := range resources {
		tributes = append(tributes, &amqp.EncyclopediaAlmanaxResourceAnswer_Tribute{
			ItemName: itemName,
			Quantity: quantity,
		})
	}

	return &amqp.EncyclopediaAlmanaxResourceAnswer{
		Tributes: tributes,
		Duration: dayDuration,
		Source:   constants.GetDofusDudeSource(),
	}
}
