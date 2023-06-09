package encyclopedias

import (
	"context"
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New(broker amqp.MessageBroker) (*Impl, error) {
	config := dodugo.NewConfiguration()
	config.UserAgent = constants.UserAgent
	apiClient := dodugo.NewAPIClient(config)

	return &Impl{
		dofusdudeClient: apiClient,
		broker:          broker,
		httpTimeout:     time.Duration(viper.GetInt(constants.DofusDudeTimeout)) * time.Second,
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

func (service *Impl) consume(_ context.Context,
	message *amqp.RabbitMQMessage, correlationID string) {
	//exhaustive:ignore Don't need to be exhaustive here since they will be handled by default case
	switch message.Type {
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_REQUEST:
		service.almanaxRequest(message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_LIST_REQUEST:
		service.itemListRequest(message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_REQUEST:
		service.itemRequest(message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_LIST_REQUEST:
		service.setListRequest(message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_REQUEST:
		service.setRequest(message, correlationID)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Type not recognized, request ignored")
	}
}
