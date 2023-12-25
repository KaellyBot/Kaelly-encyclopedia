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
