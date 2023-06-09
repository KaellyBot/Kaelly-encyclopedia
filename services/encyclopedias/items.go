package encyclopedias

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) itemListRequest(message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaItemListRequest
	if !isValidItemListRequest(request) {
		service.publishItemListAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get item list encyclopedia request received")

	dodugoItems, err := service.GetItemsAllSearch(request.Query, mappers.MapLanguage(message.Language))
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

func (service *Impl) itemRequest(message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaItemRequest
	if !isValidItemRequest(request) {
		service.publishItemAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get item encyclopedia request received")

	// TODO

	service.publishItemAnswerFailed(correlationID, message.Language)
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

func (service *Impl) publishItemAnswerSuccess(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                   amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_ANSWER,
		Status:                 amqp.RabbitMQMessage_SUCCESS,
		Language:               language,
		EncyclopediaItemAnswer: &amqp.EncyclopediaItemAnswer{
			// TODO
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}
