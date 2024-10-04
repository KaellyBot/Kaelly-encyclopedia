package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func (service *Impl) almanaxRequest(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaAlmanaxRequest
	lg := mappers.MapLanguage(message.Language)
	if !isValidAlmanaxRequest(request) {
		service.publishAlmanaxAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get almanax encyclopedia request received")

	almanax, err := service.sourceService.GetAlmanaxByDate(ctx, request.Date.AsTime(), lg)
	if err != nil {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Str(constants.LogDate, request.Date.String()).
			Msgf("Error while handling encyclopedia almanax date, returning failed request")
		service.publishAlmanaxAnswerFailed(correlationID, message.Language)
		return
	}

	service.publishAlmanaxAnswerSuccess(correlationID, mappers.MapAlmanax(almanax), message.Language)
}

func (service *Impl) almanaxEffectRequest(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaAlmanaxEffectRequest
	lg := mappers.MapLanguage(message.Language)
	if !isValidAlmanaxEffectRequest(request) {
		service.publishAlmanaxEffectAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get almanax effect encyclopedia request received")

	values, err := service.sourceService.SearchAlmanaxEffects(ctx, request.Query, lg)
	if err != nil {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Err(err).
			Str(constants.LogQueryID, request.Query).
			Msgf("Error while handling encyclopedia almanax effect" +
				" and searching for almanax effect list, returning failed request")
		service.publishAlmanaxEffectAnswerFailed(correlationID, message.Language)
		return
	}

	if len(values) == 0 {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Err(sources.ErrNotFound).
			Str(constants.LogQueryID, request.Query).
			Msgf("Error while handling encyclopedia almanax effect" +
				" and searching for almanax effect list, returning failed request")
		service.publishAlmanaxEffectAnswerFailed(correlationID, message.Language)
		return
	}

	// We trust the omnisearch by taking the first one in the list
	effect := values[0]
	almanax, err := service.sourceService.GetAlmanaxByEffect(ctx, *effect.Id, lg)
	if err != nil {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Err(err).
			Str(constants.LogQueryID, request.Query).
			Msgf("Error while handling encyclopedia almanax effect, returning failed request")
		service.publishAlmanaxEffectAnswerFailed(correlationID, message.Language)
		return
	}

	service.publishAlmanaxEffectAnswerSuccess(correlationID, mappers.MapAlmanax(almanax), message.Language)
}

func (service *Impl) almanaxResourceRequest(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaAlmanaxResourceRequest
	lg := mappers.MapLanguage(message.Language)
	if !isValidAlmanaxResourceRequest(request) {
		service.publishAlmanaxResourceAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get almanax resources encyclopedia request received")

	almanax, err := service.sourceService.GetAlmanaxByRange(ctx, request.Duration, lg)
	if err != nil {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Int32(constants.LogDuration, request.Duration).
			Msgf("Error while handling encyclopedia almanax resources, returning failed request")
		service.publishAlmanaxResourceAnswerFailed(correlationID, message.Language)
		return
	}

	answer := mappers.MapAlmanaxResource(almanax, request.Duration)
	service.publishAlmanaxResourceAnswerSuccess(correlationID, answer, message.Language)
}

func isValidAlmanaxRequest(request *amqp.EncyclopediaAlmanaxRequest) bool {
	return request != nil
}

func isValidAlmanaxEffectRequest(request *amqp.EncyclopediaAlmanaxEffectRequest) bool {
	return request != nil
}

func isValidAlmanaxResourceRequest(request *amqp.EncyclopediaAlmanaxResourceRequest) bool {
	return request != nil
}

func (service *Impl) publishAlmanaxAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishAlmanaxAnswerSuccess(correlationID string, almanax *amqp.Almanax, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		EncyclopediaAlmanaxAnswer: &amqp.EncyclopediaAlmanaxAnswer{
			Almanax: almanax,
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishAlmanaxEffectAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_EFFECT_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishAlmanaxEffectAnswerSuccess(correlationID string, almanax *amqp.Almanax,
	language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_EFFECT_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		EncyclopediaAlmanaxEffectAnswer: &amqp.EncyclopediaAlmanaxEffectAnswer{
			Almanax: almanax,
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishAlmanaxResourceAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_RESOURCE_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishAlmanaxResourceAnswerSuccess(correlationID string,
	almanax *amqp.EncyclopediaAlmanaxResourceAnswer, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                              amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_RESOURCE_ANSWER,
		Status:                            amqp.RabbitMQMessage_SUCCESS,
		Language:                          language,
		EncyclopediaAlmanaxResourceAnswer: almanax,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}
