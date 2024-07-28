package encyclopedias

import (
	"context"
	"strconv"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func (service *Impl) itemRequest(ctx context.Context,
	message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaItemRequest
	lg := mappers.MapLanguage(message.Language)
	if !isValidItemRequest(request) {
		service.publishItemAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogQueryID, request.Query).
		Str(constants.LogQueryType, request.GetType().String()).
		Msgf("Encyclopedia Item Request received")

	funcs, found := service.getItemByFuncs[request.Type]
	if !found {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while handling encyclopedia item query type, returning failed request")
		service.publishItemAnswerFailed(correlationID, message.Language)
		return
	}

	var reply *amqp.EncyclopediaItemAnswer
	var err error
	if request.GetIsID() {
		ankamaID, errID := strconv.ParseInt(request.Query, 10, 32)
		if errID != nil {
			log.Error().Err(errID).
				Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogQueryID, request.Query).
				Str(constants.LogQueryType, request.GetType().String()).
				Msgf("Error while converting query as ankamaID, returning failed request")
			service.publishItemAnswerFailed(correlationID, message.Language)
			return
		}

		reply, err = funcs.GetItemByID(ctx, int32(ankamaID), correlationID, lg)
	} else {
		reply, err = funcs.GetItemByQuery(ctx, request.Query, correlationID, lg)
	}

	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while retrieving encyclopedia item, returning failed request")
		service.publishItemAnswerFailed(correlationID, message.Language)
		return
	}

	service.publishItemAnswerSuccess(reply, correlationID, message.Language)
}

func (service *Impl) getItemByID(ctx context.Context, id int32, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {
	return nil, errBadRequestMessage
}

func (service *Impl) getItemByQuery(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error) {

	values, err := service.sourceService.SearchAnyItems(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, sources.ErrNotFound
	}

	// We trust the omnisearch by taking the first one in the list
	item := values[0]
	itemType := service.sourceService.GetItemType(item.GetItemSubtype())
	funcs, found := service.getItemByFuncs[itemType]
	if !found {
		return nil, sources.ErrNotFound
	}

	resp, err := funcs.GetItemByID(ctx, item.GetAnkamaId(), correlationID, lg)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func isValidItemRequest(request *amqp.EncyclopediaItemRequest) bool {
	return request != nil
}
