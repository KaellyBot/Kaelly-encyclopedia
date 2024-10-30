package encyclopedias

import (
	"context"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/mappers"
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

	response := mappers.MapAlmanax(almanax, service.sourceService)
	service.publishAlmanaxAnswerSuccess(correlationID, response, message.Language)
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

	effect, errEffect := service.getEffectFromRequest(ctx, request, lg)
	if errEffect != nil {
		log.Error().Str(constants.LogCorrelationID, correlationID).
			Err(errEffect).
			Str(constants.LogQueryID, request.Query).
			Str(constants.LogDate, request.Date.String()).
			Msgf("Error while handling encyclopedia almanax effect request" +
				" and searching for accurate almanax effect, returning failed request")
		service.publishAlmanaxEffectAnswerFailed(correlationID, message.Language)
		return
	}

	offset := request.GetOffset()
	adjustedSize := offset + request.GetSize()
	dodugoAlmanaxes := make([]*dodugo.AlmanaxEntry, 0)
	almanaxDates := service.almanaxService.GetDatesByAlmanaxEffect(*effect.Id)
	for i := offset; i < adjustedSize && i < int32(len(almanaxDates)); i++ {
		dodugoAlmanax, err := service.sourceService.GetAlmanaxByDate(ctx, almanaxDates[i], lg)
		if err != nil {
			log.Error().Str(constants.LogCorrelationID, correlationID).
				Str(constants.LogDate, almanaxDates[i].String()).
				Msgf("Error while handling encyclopedia almanax date, returning failed request")
			service.publishAlmanaxEffectAnswerFailed(correlationID, message.Language)
			return
		}

		dodugoAlmanaxes = append(dodugoAlmanaxes, dodugoAlmanax)
	}

	almanaxEffects := mappers.MapAlmanaxEffects(request, effect.GetName(), dodugoAlmanaxes,
		len(almanaxDates), service.sourceService)
	service.publishAlmanaxEffectAnswerSuccess(correlationID, almanaxEffects, message.Language)
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

	answer := mappers.MapAlmanaxResource(almanax, request.Duration, service.sourceService)
	service.publishAlmanaxResourceAnswerSuccess(correlationID, answer, message.Language)
}

func (service *Impl) getEffectFromRequest(ctx context.Context, request *amqp.EncyclopediaAlmanaxEffectRequest,
	lg string) (*dodugo.GetMetaAlmanaxBonuses200ResponseInner, error) {
	switch request.Type {
	case amqp.EncyclopediaAlmanaxEffectRequest_QUERY:
		values, errSearch := service.sourceService.SearchAlmanaxEffects(ctx, request.Query, lg)
		if errSearch != nil {
			return nil, errSearch
		}

		if len(values) == 0 {
			return nil, errResponseRequestEmpty
		}

		// We trust the omnisearch by taking the first one in the list
		return &values[0], nil
	case amqp.EncyclopediaAlmanaxEffectRequest_DATE:
		dodugoAlmanax, errGet := service.sourceService.
			GetAlmanaxByDate(ctx, request.GetDate().AsTime(), lg)
		if errGet != nil {
			return nil, errGet
		}
		return dodugoAlmanax.Bonus.Type, nil
	default:
		return nil, errUnknownQuery
	}
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
			Source:  constants.GetDofusDudeSource(),
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

func (service *Impl) publishAlmanaxEffectAnswerSuccess(correlationID string,
	almanaxEffects *amqp.EncyclopediaAlmanaxEffectAnswer, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                            amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_EFFECT_ANSWER,
		Status:                          amqp.RabbitMQMessage_SUCCESS,
		Language:                        language,
		EncyclopediaAlmanaxEffectAnswer: almanaxEffects,
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
