package application

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	equipmentRepo "github.com/kaellybot/kaelly-encyclopedia/repositories/equipments"
	setRepo "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/encyclopedias"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
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

	scheduler, errScheduler := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if errScheduler != nil {
		return nil, errScheduler
	}

	// Repositories
	equipmentRepo := equipmentRepo.New(db)
	setRepo := setRepo.New(db)

	// services
	storeService := stores.New()
	equipmentService, err := equipments.New(equipmentRepo)
	if err != nil {
		return nil, err
	}

	setService, err := sets.New(scheduler, setRepo)
	if err != nil {
		return nil, err
	}

	sourceService, err := sources.New(storeService)
	if err != nil {
		return nil, err
	}

	encyclopediaService := encyclopedias.New(broker, sourceService,
		equipmentService, setService)

	return &Impl{
		db:                  db,
		broker:              broker,
		scheduler:           scheduler,
		encyclopediaService: encyclopediaService,
	}, nil
}

func (app *Impl) Run() error {
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
