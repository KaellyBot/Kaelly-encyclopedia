package mappers

import amqp "github.com/kaellybot/kaelly-amqp"

func MapList(list *amqp.EncyclopediaListAnswer, language amqp.Language) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:                   amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_ANSWER,
		Status:                 amqp.RabbitMQMessage_SUCCESS,
		Language:               language,
		EncyclopediaListAnswer: list,
	}
}
