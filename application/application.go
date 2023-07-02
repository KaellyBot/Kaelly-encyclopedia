package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	equipmentRepo "github.com/kaellybot/kaelly-encyclopedia/repositories/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/encyclopedias"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	db, err := databases.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("DB instantiation failed, shutting down.")
	}

	broker, err := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		[]amqp.Binding{encyclopedias.GetBinding()})
	if err != nil {
		return nil, err
	}

	// Repositories
	equipmentRepo := equipmentRepo.New(db)

	// services
	storeService := stores.New()
	equipmentService, err := equipments.New(equipmentRepo)
	if err != nil {
		return nil, err
	}

	encyclopediaService, err := encyclopedias.New(broker, storeService, equipmentService)
	if err != nil {
		return nil, err
	}

	return &Impl{
		db:                  db,
		broker:              broker,
		encyclopediaService: encyclopediaService,
	}, nil
}

func (app *Impl) Run() error {
	return app.encyclopediaService.Consume()
}

func (app *Impl) Shutdown() {
	app.db.Shutdown()
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
