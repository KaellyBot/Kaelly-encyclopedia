package sources

import (
	"github.com/dofusdude/dodugo"
	"github.com/go-co-op/gocron/v2"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/repositories/games"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
	"github.com/spf13/viper"
)

func New(scheduler gocron.Scheduler, storeService stores.Service,
	gameRepo games.Repository) (*Impl, error) {
	config := dodugo.NewConfiguration()
	config.UserAgent = constants.UserAgent
	apiClient := dodugo.NewAPIClient(config)

	service := Impl{
		eventHandlers:   make([]GameEventHandler, 0),
		dofusDudeClient: apiClient,
		storeService:    storeService,
		gameRepo:        gameRepo,
		httpTimeout:     viper.GetDuration(constants.DofusDudeTimeout),
		itemTypes: map[string]amqp.ItemType{
			"consumables":     amqp.ItemType_CONSUMABLE_TYPE,
			"equipment":       amqp.ItemType_EQUIPMENT_TYPE,
			"items-cosmetics": amqp.ItemType_COSMETIC_TYPE,
			"items-equipment": amqp.ItemType_EQUIPMENT_TYPE,
			"mounts":          amqp.ItemType_MOUNT_TYPE,
			"quest":           amqp.ItemType_QUEST_ITEM_TYPE,
			"resources":       amqp.ItemType_RESOURCE_TYPE,
			"sets":            amqp.ItemType_SET_TYPE,
		},
	}

	_, errJob := scheduler.NewJob(
		gocron.CronJob(viper.GetString(constants.UpdateSetCronTab), true),
		gocron.NewTask(func() { service.checkGameVersion() }),
		gocron.WithName("Check game version"),
	)
	if errJob != nil {
		return nil, errJob
	}

	return &service, nil
}
