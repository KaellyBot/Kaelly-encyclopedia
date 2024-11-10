package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) listRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.EncyclopediaListRequest
	if !isValidListRequest(request) {
		service.replyWithFailedAnswer(ctx, amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_ANSWER,
			message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogQueryID, request.Query).
		Str(constants.LogQueryType, request.GetType().String()).
		Msgf("Encyclopedia List Request received")

	getListFunc, found := service.getListByFunc[request.Type]
	if !found {
		log.Error().Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while handling encyclopedia list query type, returning failed request")
		service.replyWithFailedAnswer(ctx, amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_ANSWER,
			message.Language)
		return
	}

	list, err := getListFunc(ctx, request.Query, ctx.CorrelationID,
		mappers.MapLanguage(message.Language))
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogQueryType, request.GetType().String()).
			Msgf("Error while retrieving encyclopedia list, returning failed request")
		service.replyWithFailedAnswer(ctx, amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_ANSWER,
			message.Language)
		return
	}

	response := mappers.MapList(list, message.Language)
	service.replyWithSuceededAnswer(ctx, response)
}

func (service *Impl) getItemList(ctx context.Context, query, _,
	lg string) (*amqp.EncyclopediaListAnswer, error) {
	dodugoItems, err := service.sourceService.SearchAnyItems(ctx, query, lg)
	if err != nil {
		return nil, err
	}

	return mappers.MapItemList(dodugoItems), nil
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

func isValidListRequest(request *amqp.EncyclopediaListRequest) bool {
	return request != nil
}
