package application

import (
	"github.com/go-co-op/gocron/v2"
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
	scheduler           gocron.Scheduler
	encyclopediaService encyclopedias.Service
}
