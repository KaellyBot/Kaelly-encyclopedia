package application

import (
	"github.com/go-co-op/gocron/v2"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/services/encyclopedias"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
	"github.com/kaellybot/kaelly-encyclopedia/utils/insights"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	broker              amqp.MessageBroker
	scheduler           gocron.Scheduler
	db                  databases.MySQLConnection
	probes              insights.Probes
	prom                insights.PrometheusMetrics
	encyclopediaService encyclopedias.Service
}
