package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) itemListRequest(ctx context.Context,
	message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaItemListRequest
	if !isValidItemListRequest(request) {
		service.publishItemListAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogQueryID, request.Query).
		Str(constants.LogQueryType, request.GetType().String()).
		Msgf("Encyclopedia List Request received")

	getItemListFunc, found := service.getItemListByFunc[request.Type]
	if !found {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while handling encyclopedia list query type, returning failed request")
		service.publishItemListAnswerFailed(correlationID, message.Language)
		return
	}

	reply, err := getItemListFunc(ctx, request.Query, correlationID,
		mappers.MapLanguage(message.Language))
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while retrieving encyclopedia list, returning failed request")
		service.publishItemListAnswerFailed(correlationID, message.Language)
		return
	}

	service.publishItemListAnswerSuccess(reply, correlationID, message.Language)
}

func (service *Impl) getItemList(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemListAnswer, error) {
	dodugoItems, err := service.sourceService.SearchAnyItems(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapItemList(dodugoItems), nil
}

func (service *Impl) getSetList(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemListAnswer, error) {
	dodugoSets, err := service.sourceService.SearchSets(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapSetList(dodugoSets), nil
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

func (service *Impl) publishItemListAnswerSuccess(answer *amqp.EncyclopediaItemListAnswer,
	correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                       amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_LIST_ANSWER,
		Status:                     amqp.RabbitMQMessage_SUCCESS,
		Language:                   language,
		EncyclopediaItemListAnswer: answer,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func isValidItemListRequest(request *amqp.EncyclopediaItemListRequest) bool {
	return request != nil
}
