package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func (service *Impl) almanaxRequest(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaAlmanaxRequest
	if !isValidAlmanaxRequest(request) {
		service.publishAlmanaxAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get almanax encyclopedia request received")

	// TODO

	service.publishAlmanaxAnswerFailed(correlationID, message.Language)
}

func isValidAlmanaxRequest(request *amqp.EncyclopediaAlmanaxRequest) bool {
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

func (service *Impl) publishAlmanaxAnswerSuccess(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                      amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_ANSWER,
		Status:                    amqp.RabbitMQMessage_SUCCESS,
		Language:                  language,
		EncyclopediaAlmanaxAnswer: &amqp.EncyclopediaAlmanaxAnswer{
			// TODO
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}
