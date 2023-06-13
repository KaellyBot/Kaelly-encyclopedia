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

func (service *Impl) GetItemsAllSearch(ctx context.Context, query,
	language string) ([]dodugo.ItemsListEntryTyped, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var items []dodugo.ItemsListEntryTyped
	key := buildDofusDudeKey(item, query, language)
	if !service.getListFromCache(ctx, key, &items) {
		resp, r, err := service.dofusDudeClient.AllItemsApi.
			GetItemsAllSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()
		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		service.putListToCache(ctx, key, resp)
		items = resp
	}

	return items, nil
}

func (service *Impl) GetSetsSearch(ctx context.Context, query,
	language string) ([]dodugo.SetListEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var sets []dodugo.SetListEntry
	key := buildDofusDudeKey(set, query, language)
	if !service.getListFromCache(ctx, key, &sets) {
		resp, r, err := service.dofusDudeClient.SetsApi.
			GetSetsSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()
		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		service.putListToCache(ctx, key, resp)
		sets = resp
	}

	return sets, nil
}

func (service *Impl) getListFromCache(ctx context.Context, key string, value any) bool {
	err := service.storeService.Get(ctx, key, value)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			log.Info().
				Str(constants.LogKey, key).
				Msgf("Cannot find elements in cache, calling the API...")
		} else {
			log.Error().Err(err).
				Str(constants.LogKey, key).
				Msgf("Error while requesting elements in cache, calling the API instead...")
		}
	}

	return err != nil
}

func (service *Impl) putListToCache(ctx context.Context, key string, value any) {
	err := service.storeService.Set(ctx, key, value)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogKey, key).
			Msgf("Error while putting elments in cache, no issue to retrieve elments anyway...")
	}
}

func buildDofusDudeKey(objType objectType, query, language string) string {
	return fmt.Sprintf("%v/%v?query=%v&lg=%v", constants.GetEncyclopediasSource().Name,
		objType, query, language)
}
