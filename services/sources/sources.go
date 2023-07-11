package sources

import (
	"github.com/dofusdude/dodugo"
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
	}, nil
}
