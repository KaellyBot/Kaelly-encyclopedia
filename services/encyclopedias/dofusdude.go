package encyclopedias

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dofusdude/dodugo"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
)

func (service *Impl) searchItems(ctx context.Context, query,
	language string) ([]dodugo.ItemsListEntryTyped, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var items []dodugo.ItemsListEntryTyped
	key := buildListKey(item, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &items) {
		resp, r, err := service.dofusDudeClient.AllItemsAPI.
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

func (service *Impl) getItemByQuery(ctx context.Context, query, language string,
) (*dodugo.Weapon, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	values, err := service.searchItems(ctx, query, language)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, errNotFound
	}

	// We trust the omnisearch by taking the first one in the list
	resp, err := service.getItemByID(ctx, values[0].GetAnkamaId(), language)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Impl) getItemByID(ctx context.Context, itemID int32, language string,
) (*dodugo.Weapon, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoItem *dodugo.Weapon
	key := buildKey(item, fmt.Sprintf("%v", itemID), language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoItem) {
		resp, r, err := service.dofusDudeClient.EquipmentAPI.
			GetItemsEquipmentSingle(ctx, language, itemID, constants.DofusDudeGame).Execute()
		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoItem = resp
	}

	return dodugoItem, nil
}

func (service *Impl) searchSets(ctx context.Context, query,
	language string) ([]dodugo.SetListEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var sets []dodugo.SetListEntry
	key := buildListKey(set, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &sets) {
		resp, r, err := service.dofusDudeClient.SetsAPI.
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

func (service *Impl) getSetByQuery(ctx context.Context, query, language string,
) (*dodugo.EquipmentSet, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	values, err := service.searchSets(ctx, query, language)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, errNotFound
	}

	// We trust the omnisearch by taking the first one in the list
	resp, err := service.getSetByID(ctx, values[0].GetAnkamaId(), language)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Impl) getSetByID(ctx context.Context, setID int32, language string,
) (*dodugo.EquipmentSet, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoSet *dodugo.EquipmentSet
	key := buildKey(set, fmt.Sprintf("%v", setID), language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoSet) {
		resp, r, err := service.dofusDudeClient.SetsAPI.
			GetSetsSingle(ctx, language, setID, constants.DofusDudeGame).Execute()
		if err != nil && r.StatusCode != http.StatusNotFound {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoSet = resp
	}

	return dodugoSet, nil
}
