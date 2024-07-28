package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) listRequest(ctx context.Context,
	message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaListRequest
	if !isValidListRequest(request) {
		service.publishListAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogQueryID, request.Query).
		Str(constants.LogQueryType, request.GetType().String()).
		Msgf("Encyclopedia List Request received")

	getListFunc, found := service.getListByFunc[request.Type]
	if !found {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while handling encyclopedia list query type, returning failed request")
		service.publishListAnswerFailed(correlationID, message.Language)
		return
	}

	reply, err := getListFunc(ctx, request.Query, correlationID,
		mappers.MapLanguage(message.Language))
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while retrieving encyclopedia list, returning failed request")
		service.publishListAnswerFailed(correlationID, message.Language)
		return
	}

	service.publishListAnswerSuccess(reply, correlationID, message.Language)
}

func (service *Impl) getItemList(ctx context.Context, query, _,
	lg string) (*amqp.EncyclopediaListAnswer, error) {
	dodugoItems, err := service.sourceService.SearchEquipments(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapEquipmentList(dodugoItems), nil
}

func (service *Impl) getSetList(ctx context.Context, query, _,
	lg string) (*amqp.EncyclopediaListAnswer, error) {
	dodugoSets, err := service.sourceService.SearchSets(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapSetList(dodugoSets), nil
}

func (service *Impl) getAlmanaxEffectList(ctx context.Context, query, _,
	lg string) (*amqp.EncyclopediaListAnswer, error) {
	dodugoAlmanaxEffects, err := service.sourceService.SearchAlmanaxEffects(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapAlmanaxEffectList(dodugoAlmanaxEffects), nil
}

func (service *Impl) publishListAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishListAnswerSuccess(answer *amqp.EncyclopediaListAnswer,
	correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                   amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_ANSWER,
		Status:                 amqp.RabbitMQMessage_SUCCESS,
		Language:               language,
		EncyclopediaListAnswer: answer,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func isValidListRequest(request *amqp.EncyclopediaListRequest) bool {
	return request != nil
}
