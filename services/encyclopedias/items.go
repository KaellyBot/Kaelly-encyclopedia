package encyclopedias

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) itemListRequest(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaItemListRequest
	if !isValidItemListRequest(request) {
		service.publishItemListAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogQueryID, request.Query).
		Msgf("Get item list encyclopedia request received")

	dodugoItems, err := service.searchItems(ctx, request.Query, mappers.MapLanguage(message.Language))
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Msgf("Error while calling DofusDude, returning failed request")
		service.publishItemListAnswerFailed(correlationID, message.Language)
		return
	}

	items := mappers.MapItemList(dodugoItems)
	service.publishItemListAnswerSuccess(items, correlationID, message.Language)
}

func (service *Impl) itemRequest(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaItemRequest
	if !isValidItemRequest(request) {
		service.publishItemAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get item encyclopedia request received")

	var item *dodugo.Weapon
	var err error
	lg := mappers.MapLanguage(message.Language)
	if request.GetIsID() {
		ankamaID, err := strconv.Atoi(request.Query)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogQueryID, request.Query).
				Msgf("Error while converting query as AnkamaID, returning failed request")
			service.publishItemAnswerFailed(correlationID, message.Language)
			return
		}

		item, err = service.getItemByID(ctx, int32(ankamaID), lg)
	} else {
		item, err = service.getItemByQuery(ctx, request.Query, lg)
	}
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Msgf("Error while calling DofusDude, returning failed request")
		service.publishItemAnswerFailed(correlationID, message.Language)
		return
	}

	ingredients := make(map[int32]*dodugo.Weapon)
	for _, ingredient := range item.Recipe {
		itemID := ingredient.GetItemAnkamaId()
		// TODO getItemByID depends of the searched type!
		item, errItem := service.getItemByID(ctx, itemID, lg)
		if errItem != nil {
			log.Error().Err(errItem).
				Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogQueryID, request.Query).
				Str(constants.LogAnkamaID, fmt.Sprintf("%v", itemID)).
				Msgf("Error while retrieving item with DofusDude, continuing without it")
		} else {
			ingredients[itemID] = item
		}
	}

	answer := mappers.MapItem(item, ingredients)
	service.publishItemAnswerSuccess(answer, correlationID, message.Language)
}

func isValidItemListRequest(request *amqp.EncyclopediaItemListRequest) bool {
	return request != nil
}

func isValidItemRequest(request *amqp.EncyclopediaItemRequest) bool {
	return request != nil
}

func (service *Impl) publishItemListAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_LIST_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishItemListAnswerSuccess(items []*amqp.EncyclopediaItemListAnswer_Item,
	correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_LIST_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		EncyclopediaItemListAnswer: &amqp.EncyclopediaItemListAnswer{
			Items: items,
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishItemAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishItemAnswerSuccess(answer *amqp.EncyclopediaItemAnswer,
	correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                   amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_ANSWER,
		Status:                 amqp.RabbitMQMessage_SUCCESS,
		Language:               language,
		EncyclopediaItemAnswer: answer,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}
