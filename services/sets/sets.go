package sets

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New(scheduler gocron.Scheduler,
	repository repository.Repository) (*Impl, error) {
	service := Impl{
		sets:       make(map[int32]entities.Set),
		repository: repository,
	}

	errDB := service.loadSetFromDB()
	if errDB != nil {
		return nil, errDB
	}

	_, errJob := scheduler.NewJob(
		gocron.CronJob(viper.GetString(constants.UpdateSetCronTab), true),
		gocron.NewTask(func() { service.updateSets() }),
		gocron.WithName("Set icons build"),
	)
	if errJob != nil {
		return nil, errJob
	}

	return &service, nil
}

func (service *Impl) GetSetByDofusDude(id int32) (entities.Set, bool) {
	item, found := service.sets[id]
	return item, found
}

func (service *Impl) loadSetFromDB() error {
	sets, err := service.repository.GetSets()
	if err != nil {
		return err
	}

	log.Info().
		Int(constants.LogEntityCount, len(sets)).
		Msgf("Sets loaded")

	for _, set := range sets {
		service.sets[set.DofusDudeID] = set
	}

	return nil
}

func (service *Impl) updateSets() {
	// TODO
}
