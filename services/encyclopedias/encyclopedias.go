package encyclopedias

import (
	"context"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New(broker amqp.MessageBroker) (*Impl, error) {
	// TODO init dodugo

	return &Impl{
		broker:      broker,
		httpTimeout: time.Duration(viper.GetInt(constants.DofusDudeTimeout)) * time.Second,
	}, nil
}

func GetBinding() amqp.Binding {
	return amqp.Binding{
		Exchange:   amqp.ExchangeRequest,
		RoutingKey: requestsRoutingkey,
		Queue:      requestQueueName,
	}
}

func (service *Impl) Consume() error {
	log.Info().Msgf("Consuming encyclopedia requests...")
	return service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(ctx context.Context, message *amqp.RabbitMQMessage, correlationID string) {
	if !isValidEncyclopediaRequest(message) {
		service.publishEncyclopediaAnswerFailed(correlationID, message.Language)
		return
	}

	// TODO

	service.publishEncyclopediaAnswerSuccess(correlationID, message.Language)
}

func isValidEncyclopediaRequest(message *amqp.RabbitMQMessage) bool {
	// TODO
	return message.Type == amqp.RabbitMQMessage_PORTAL_POSITION_REQUEST && message.GetPortalPositionRequest() != nil
}

func (service *Impl) publishEncyclopediaAnswerFailed(correlationID string, language amqp.Language) {
	// TODO
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishEncyclopediaAnswerSuccess(correlationID string, language amqp.Language) {
	// TODO
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_PORTAL_POSITION_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}
