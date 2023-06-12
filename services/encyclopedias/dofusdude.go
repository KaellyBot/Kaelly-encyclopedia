package encyclopedias

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dofusdude/dodugo"
	"github.com/go-redis/cache/v8"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetItemsAllSearch(ctx context.Context, query, language string) ([]dodugo.ItemsListEntryTyped, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var items []dodugo.ItemsListEntryTyped
	key := buildDofusDudeKey(item, query, language)
	err := service.storeService.Get(ctx, key, &items)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			log.Info().
				Str(constants.LogKey, key).
				Msgf("Cannot find items in cache, calling the API...")
		} else {
			log.Error().Err(err).
				Str(constants.LogKey, key).
				Msgf("Error while requesting items in cache, calling the API instead...")
		}

		resp, r, err := service.dofusDudeClient.AllItemsApi.
			GetItemsAllSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()

		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		items = resp

		err = service.storeService.Set(ctx, key, items)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogKey, key).
				Msgf("Error while putting items in cache, no issue to retrieve items anyway...")
		}
	}

	return items, nil
}

func buildDofusDudeKey(objType objectType, query, language string) string {
	return fmt.Sprintf("%v/%v?query=%v&lg=%v", constants.GetEncyclopediasSource().Name,
		objType, query, language)
}
