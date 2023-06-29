package encyclopedias

import (
	"context"
	"net/http"

	"github.com/dofusdude/dodugo"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
)

func (service *Impl) GetItemsAllSearch(ctx context.Context, query,
	language string) ([]dodugo.ItemsListEntryTyped, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var items []dodugo.ItemsListEntryTyped
	key := buildListKey(item, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &items) {
		resp, r, err := service.dofusDudeClient.AllItemsApi.
			GetItemsAllSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()
		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		items = resp
	}

	return items, nil
}

func (service *Impl) GetSetsSearch(ctx context.Context, query,
	language string) ([]dodugo.SetListEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var sets []dodugo.SetListEntry
	key := buildListKey(set, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &sets) {
		resp, r, err := service.dofusDudeClient.SetsApi.
			GetSetsSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()
		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		sets = resp
	}

	return sets, nil
}

func (service *Impl) GetSet(ctx context.Context, query,
	language string) (*dodugo.EquipmentSet, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoSet *dodugo.EquipmentSet
	key := buildKey(set, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoSet) {

		// TODO determine Ankama ID
		values, err := service.GetSetsSearch(ctx, query, language)
		if err != nil {
			return nil, err
		}
		var ankamaID int32 = *values[0].AnkamaId
		// --------------------------

		resp, r, err := service.dofusDudeClient.SetsApi.
			GetSetsSingle(ctx, language, ankamaID, constants.DofusDudeGame).Execute()
		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoSet = resp
	}

	return dodugoSet, nil
}
