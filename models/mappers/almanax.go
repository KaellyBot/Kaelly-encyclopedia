package mappers

import (
	"fmt"
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapAlmanaxEffects(request *amqp.EncyclopediaAlmanaxEffectRequest, effectName string,
	dodugoAlmanaxes []*dodugo.AlmanaxEntry, total int, sourceService sources.Service,
	language amqp.Language) *amqp.RabbitMQMessage {
	almanaxes := make([]*amqp.Almanax, 0)
	for _, dodugoAlmanax := range dodugoAlmanaxes {
		almanax := mapAlmanax(dodugoAlmanax, sourceService)
		if almanax == nil {
			return nil
		}

		almanaxes = append(almanaxes, almanax)
	}

	page := request.GetOffset() / request.GetSize()
	pages := int32(total) / request.GetSize()
	if int32(total)%request.GetSize() != 0 {
		pages++
	}

	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_EFFECT_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		EncyclopediaAlmanaxEffectAnswer: &amqp.EncyclopediaAlmanaxEffectAnswer{
			Query:     effectName,
			Almanaxes: almanaxes,
			Page:      page,
			Pages:     pages,
			Total:     int32(total),
			Source:    constants.GetDofusDudeSource(),
		},
	}
}

func MapAlmanax(dodugoAlmanax *dodugo.AlmanaxEntry, sourceService sources.Service,
	language amqp.Language) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		EncyclopediaAlmanaxAnswer: &amqp.EncyclopediaAlmanaxAnswer{
			Almanax: mapAlmanax(dodugoAlmanax, sourceService),
			Source:  constants.GetDofusDudeSource(),
		},
	}
}

func mapAlmanax(dodugoAlmanax *dodugo.AlmanaxEntry, sourceService sources.Service,
) *amqp.Almanax {
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

	itemType := sourceService.GetItemType(dodugoAlmanax.Tribute.Item.GetSubtype())

	return &amqp.Almanax{
		Bonus: *dodugoAlmanax.Bonus.Description,
		Tribute: &amqp.Almanax_Tribute{
			Item: &amqp.Almanax_Tribute_Item{
				Name: *dodugoAlmanax.Tribute.Item.Name,
				Icon: icon,
				Type: itemType,
			},
			Quantity: *dodugoAlmanax.Tribute.Quantity,
		},
		Reward: int64(dodugoAlmanax.GetRewardKamas()),
		Date:   timestamppb.New(date.UTC()),
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
	sourceService sources.Service, language amqp.Language) *amqp.RabbitMQMessage {
	quantityPerResource := make(map[string]int32, 0)
	for _, almanax := range dodugoAlmanax {
		itemName := *almanax.Tribute.GetItem().Name
		quantity, found := quantityPerResource[itemName]
		if !found {
			quantity = 0
		}

		quantityPerResource[itemName] = quantity + almanax.Tribute.GetQuantity()
	}

	tributes := make([]*amqp.EncyclopediaAlmanaxResourceAnswer_Tribute, 0)
	for _, almanax := range dodugoAlmanax {
		itemName := *almanax.Tribute.GetItem().Name
		tributes = append(tributes, &amqp.EncyclopediaAlmanaxResourceAnswer_Tribute{
			ItemName: itemName,
			ItemType: sourceService.GetItemType(almanax.Tribute.Item.GetSubtype()),
			Quantity: quantityPerResource[itemName],
		})
	}

	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_RESOURCE_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		EncyclopediaAlmanaxResourceAnswer: &amqp.EncyclopediaAlmanaxResourceAnswer{
			Tributes: tributes,
			Duration: dayDuration,
			Source:   constants.GetDofusDudeSource(),
		},
	}
}
