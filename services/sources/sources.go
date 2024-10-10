package sources

import (
	"github.com/dofusdude/dodugo"
	"github.com/go-co-op/gocron/v2"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
	"github.com/spf13/viper"
)

func New(scheduler gocron.Scheduler,
	storeService stores.Service) (*Impl, error) {
	config := dodugo.NewConfiguration()
	config.UserAgent = constants.UserAgent
	apiClient := dodugo.NewAPIClient(config)

	service := Impl{
		eventHandlers:   make([]GameEventHandler, 0),
		dofusDudeClient: apiClient,
		storeService:    storeService,
		httpTimeout:     viper.GetDuration(constants.DofusDudeTimeout),
		itemTypes: map[string]amqp.ItemType{
			"items-equipment": amqp.ItemType_EQUIPMENT,
			"mounts":          amqp.ItemType_MOUNT_TYPE,
			"sets":            amqp.ItemType_SET,
		},
		ingredientTypes: map[string]amqp.IngredientType{
			"consumables": amqp.IngredientType_CONSUMABLE,
			"equipment":   amqp.IngredientType_EQUIPMENT_INGREDIENT,
			"quest":       amqp.IngredientType_QUEST_ITEM,
			"resources":   amqp.IngredientType_RESOURCE,
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
