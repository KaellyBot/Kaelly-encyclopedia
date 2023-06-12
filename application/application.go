package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/encyclopedias"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	broker, err := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		[]amqp.Binding{encyclopedias.GetBinding()})
	if err != nil {
		return nil, err
	}

	// services
	storeService := stores.New()

	encyclopediaService, err := encyclopedias.New(broker, storeService)
	if err != nil {
		return nil, err
	}

	return &Impl{
		encyclopediaService: encyclopediaService,
		broker:              broker,
	}, nil
}

func (app *Impl) Run() error {
	return app.encyclopediaService.Consume()
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
