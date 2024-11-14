package news

import (
	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	newsAlmanaxRoutingKey = "news.almanax"
	newsGameRoutingKey    = "news.game"
	newsSetRoutingKey     = "news.set"
)

type Service interface {
	PublishAlmanaxNews()
	PublishGameNews(gameVersion string)
	PublishSetNews(missingSetNumber, buildSetNumber int)
}

type Impl struct {
	broker amqp.MessageBroker
}
