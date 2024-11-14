package application

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	almanaxRepo "github.com/kaellybot/kaelly-encyclopedia/repositories/almanaxes"
	equipmentRepo "github.com/kaellybot/kaelly-encyclopedia/repositories/equipments"
	setRepo "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/almanaxes"
	"github.com/kaellybot/kaelly-encyclopedia/services/encyclopedias"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/news"
	"github.com/kaellybot/kaelly-encyclopedia/services/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	db, errDB := databases.New()
	if errDB != nil {
		log.Fatal().Err(errDB).Msgf("DB instantiation failed, shutting down.")
	}

	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		amqp.WithBindings(encyclopedias.GetBinding()))

	// Create scheduler with Europe/Paris timezone
	frenchLocation, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return nil, err
	}

	// Since we have winter/summer hours, UTC location cannot be used easily.
	scheduler, errScheduler := gocron.NewScheduler(gocron.WithLocation(frenchLocation))
	if errScheduler != nil {
		return nil, errScheduler
	}

	// Repositories
	almanaxRepo := almanaxRepo.New(db)
	equipmentRepo := equipmentRepo.New(db)
	setRepo := setRepo.New(db)

	// services
	storeService := stores.New()
	equipmentService, errEquipment := equipments.New(equipmentRepo)
	if errEquipment != nil {
		return nil, errEquipment
	}

	sourceService, errSource := sources.New(scheduler, storeService)
	if errSource != nil {
		return nil, errSource
	}

	newsService := news.New(broker, sourceService)
	almanaxService, errAlmanax := almanaxes.New(scheduler, frenchLocation,
		almanaxRepo, sourceService, newsService)
	if errAlmanax != nil {
		return nil, errAlmanax
	}

	setService, errSet := sets.New(setRepo, newsService, sourceService, equipmentService)
	if errSet != nil {
		return nil, errSet
	}

	encyclopediaService := encyclopedias.New(broker, sourceService,
		almanaxService, equipmentService, setService)

	return &Impl{
		db:                  db,
		broker:              broker,
		scheduler:           scheduler,
		encyclopediaService: encyclopediaService,
	}, nil
}

func (app *Impl) Run() error {
	errBroker := app.broker.Run()
	if errBroker != nil {
		return errBroker
	}

	app.scheduler.Start()
	for _, job := range app.scheduler.Jobs() {
		scheduledTime, err := job.NextRun()
		if err == nil {
			log.Info().Msgf("%v scheduled at %v", job.Name(), scheduledTime)
		}
	}

	return app.encyclopediaService.Consume()
}

func (app *Impl) Shutdown() {
	if err := app.scheduler.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Cannot shutdown scheduler, continuing...")
	}

	app.db.Shutdown()
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
