package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/services/encyclopedias"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	encyclopediaService encyclopedias.Service
	broker              amqp.MessageBroker
}
