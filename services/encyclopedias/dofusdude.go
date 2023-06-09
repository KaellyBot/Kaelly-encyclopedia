package encyclopedias

import (
	"context"

	"github.com/dofusdude/dodugo"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
)

func (service *Impl) GetItemsAllSearch(query, language string) ([]dodugo.ItemsListEntryTyped, error) {
	// TODO build context with timeout
	resp, r, err := service.dofusdudeClient.AllItemsApi.
		GetItemsAllSearch(context.Background(), language, constants.DofusDudeGame).
		Query(query).Limit(constants.DofusDudeLimit).Execute()
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return resp, nil
}
