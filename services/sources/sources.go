package sources

import (
	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
	"github.com/spf13/viper"
)

func New(storeService stores.Service) (*Impl, error) {
	config := dodugo.NewConfiguration()
	config.UserAgent = constants.UserAgent
	apiClient := dodugo.NewAPIClient(config)

	return &Impl{
		dofusDudeClient: apiClient,
		storeService:    storeService,
		httpTimeout:     viper.GetDuration(constants.DofusDudeTimeout),
		itemTypes: map[string]amqp.ItemType{
			"consumables": amqp.ItemType_CONSUMABLE,
			"cosmetics":   amqp.ItemType_COSMETIC,
			"equipment":   amqp.ItemType_EQUIPMENT,
			"mounts":      amqp.ItemType_MOUNT,
			"quest":       amqp.ItemType_QUEST_ITEM,
			"resources":   amqp.ItemType_RESOURCE,
			"sets":        amqp.ItemType_SET,
		},
	}, nil
}
