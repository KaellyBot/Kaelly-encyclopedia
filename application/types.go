package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/services/encyclopedias"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	db                  databases.MySQLConnection
	broker              amqp.MessageBroker
	encyclopediaService encyclopedias.Service
}
