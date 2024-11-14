package mappers

import amqp "github.com/kaellybot/kaelly-amqp"

func MapAlmanaxNews() *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:               amqp.RabbitMQMessage_NEWS_ALMANAX,
		Language:           amqp.Language_ANY,
		Game:               amqp.Game_DOFUS_GAME,
		NewsAlmanaxMessage: &amqp.NewsAlmanaxMessage{
			// TODO
		},
	}
}

func MapGameNews(gameVersion string) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_NEWS_GAME,
		Language: amqp.Language_ANY,
		Game:     amqp.Game_DOFUS_GAME,
		NewsGameMessage: &amqp.NewsGameMessage{
			Version: gameVersion,
		},
	}
}

func MapSetNews(missingSetNumber, buildSetNumber int) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_NEWS_SET,
		Language: amqp.Language_ANY,
		Game:     amqp.Game_DOFUS_GAME,
		NewsSetMessage: &amqp.NewsSetMessage{
			MissingSetNumber: int64(missingSetNumber),
			BuiltSetNumber:   int64(buildSetNumber),
		},
	}
}
