package encyclopedias

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func (service *Impl) setListRequest(message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaSetListRequest
	if !isValidSetListRequest(request) {
		service.publishSetListAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get set list encyclopedia request received")

	// TODO

	service.publishSetListAnswerFailed(correlationID, message.Language)
}

func (service *Impl) setRequest(message *amqp.RabbitMQMessage, correlationID string) {
	request := message.EncyclopediaSetRequest
	if !isValidSetRequest(request) {
		service.publishSetAnswerFailed(correlationID, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get set encyclopedia request received")

	// TODO

	service.publishSetAnswerFailed(correlationID, message.Language)
}

func isValidSetListRequest(request *amqp.EncyclopediaSetListRequest) bool {
	return request != nil
}

func isValidSetRequest(request *amqp.EncyclopediaSetRequest) bool {
	return request != nil
}

func (service *Impl) publishSetListAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_LIST_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishSetListAnswerSuccess(sets []*amqp.EncyclopediaSetListAnswer_Set,
	correlationID string, language amqp.Language) {	
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_LIST_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: language,
		EncyclopediaSetListAnswer: &amqp.EncyclopediaSetListAnswer{
			Sets: sets,
		},
	}
	
	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishSetAnswerFailed(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishSetAnswerSuccess(correlationID string, language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:                   amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_ANSWER,
		Status:                 amqp.RabbitMQMessage_SUCCESS,
		Language:               language,
		EncyclopediaSetAnswer: &amqp.EncyclopediaSetAnswer{
			// TODO
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}
