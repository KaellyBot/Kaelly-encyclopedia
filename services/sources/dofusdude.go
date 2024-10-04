package sources

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetItemType(itemType string) amqp.ItemType {
	amqpItemType, found := service.itemTypes[itemType]
	if !found {
		log.Warn().
			Str(constants.LogItemType, itemType).
			Msgf("Cannot find dofusDude itemType match, returning amqp.ItemType_ANY_ITEM")
		return amqp.ItemType_ANY_ITEM
	}

	return amqpItemType
}

func (service *Impl) GetIngredientType(itemType string) amqp.IngredientType {
	amqpIngredientType, found := service.ingredientTypes[itemType]
	if !found {
		log.Warn().
			Str(constants.LogItemType, itemType).
			Msgf("Cannot find dofusDude ingredientType match, returning amqp.IngredientType_ANY_INGREDIENT")
		return amqp.IngredientType_ANY_INGREDIENT
	}

	return amqpIngredientType
}

func (service *Impl) SearchAnyItems(ctx context.Context, query,
	language string) ([]dodugo.ItemsListEntryTyped, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var items []dodugo.ItemsListEntryTyped
	key := buildListKey(item, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &items) {
		resp, r, err := service.dofusDudeClient.GameAPI.
			GetItemsAllSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		items = resp
	}

	return items, nil
}

func (service *Impl) GetConsumableByID(ctx context.Context, itemID int32, language string,
) (*dodugo.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoItem *dodugo.Resource
	key := buildKey(item, fmt.Sprintf("%v", itemID), language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoItem) {
		resp, r, err := service.dofusDudeClient.ConsumablesAPI.
			GetItemsConsumablesSingle(ctx, language, itemID, constants.DofusDudeGame).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoItem = resp
	}

	return dodugoItem, nil
}

func (service *Impl) SearchEquipments(ctx context.Context, query,
	language string) ([]dodugo.ItemListEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var items []dodugo.ItemListEntry
	key := buildListKey(item, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &items) {
		resp, r, err := service.dofusDudeClient.EquipmentAPI.
			GetItemsEquipmentSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		items = resp
	}

	return items, nil
}

func (service *Impl) GetEquipmentByQuery(ctx context.Context, query, language string,
) (*dodugo.Weapon, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	values, err := service.SearchEquipments(ctx, query, language)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, ErrNotFound
	}

	// We trust the omnisearch by taking the first one in the list
	resp, err := service.GetEquipmentByID(ctx, values[0].GetAnkamaId(), language)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Impl) GetEquipmentByID(ctx context.Context, itemID int32, language string,
) (*dodugo.Weapon, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoItem *dodugo.Weapon
	key := buildKey(item, fmt.Sprintf("%v", itemID), language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoItem) {
		resp, r, err := service.dofusDudeClient.EquipmentAPI.
			GetItemsEquipmentSingle(ctx, language, itemID, constants.DofusDudeGame).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoItem = resp
	}

	return dodugoItem, nil
}

func (service *Impl) GetQuestItemByID(ctx context.Context, itemID int32, language string,
) (*dodugo.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoItem *dodugo.Resource
	key := buildKey(item, fmt.Sprintf("%v", itemID), language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoItem) {
		resp, r, err := service.dofusDudeClient.QuestItemsAPI.
			GetItemQuestSingle(ctx, language, itemID, constants.DofusDudeGame).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoItem = resp
	}

	return dodugoItem, nil
}

func (service *Impl) GetResourceByID(ctx context.Context, itemID int32, language string,
) (*dodugo.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoItem *dodugo.Resource
	key := buildKey(item, fmt.Sprintf("%v", itemID), language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoItem) {
		resp, r, err := service.dofusDudeClient.ResourcesAPI.
			GetItemsResourcesSingle(ctx, language, itemID, constants.DofusDudeGame).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoItem = resp
	}

	return dodugoItem, nil
}

func (service *Impl) SearchSets(ctx context.Context, query,
	language string) ([]dodugo.SetListEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var sets []dodugo.SetListEntry
	key := buildListKey(set, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &sets) {
		resp, r, err := service.dofusDudeClient.SetsAPI.
			GetSetsSearch(ctx, language, constants.DofusDudeGame).
			Query(query).Limit(constants.DofusDudeLimit).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		sets = resp
	}

	return sets, nil
}

func (service *Impl) GetSetByQuery(ctx context.Context, query, language string,
) (*dodugo.EquipmentSet, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	values, err := service.SearchSets(ctx, query, language)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, ErrNotFound
	}

	// We trust the omnisearch by taking the first one in the list
	resp, err := service.GetSetByID(ctx, values[0].GetAnkamaId(), language)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Impl) GetSetByID(ctx context.Context, setID int32, language string,
) (*dodugo.EquipmentSet, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoSet *dodugo.EquipmentSet
	key := buildKey(set, fmt.Sprintf("%v", setID), language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoSet) {
		resp, r, err := service.dofusDudeClient.SetsAPI.
			GetSetsSingle(ctx, language, setID, constants.DofusDudeGame).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoSet = resp
	}

	return dodugoSet, nil
}

func (service *Impl) SearchAlmanaxEffects(ctx context.Context, query,
	language string) ([]dodugo.GetMetaAlmanaxBonuses200ResponseInner, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var effects []dodugo.GetMetaAlmanaxBonuses200ResponseInner
	key := buildListKey(almanaxEffect, query, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &effects) {
		resp, r, err := service.dofusDudeClient.MetaAPI.
			GetMetaAlmanaxBonusesSearch(ctx, language).
			Query(query).
			Limit(constants.DofusDudeLimit).
			Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		effects = resp
	}

	return effects, nil
}

func (service *Impl) GetAlmanaxByDate(ctx context.Context, date time.Time, language string,
) (*dodugo.AlmanaxEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoAlmanax *dodugo.AlmanaxEntry
	dodugoAlmanaxDate := date.Format(constants.DofusDudeAlmanaxDateFormat)
	key := buildKey(almanax, dodugoAlmanaxDate, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoAlmanax) {
		resp, r, err := service.dofusDudeClient.AlmanaxAPI.
			GetAlmanaxDate(ctx, language, dodugoAlmanaxDate).Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoAlmanax = resp
	}

	if dodugoAlmanax == nil {
		log.Warn().
			Str(constants.LogDate, dodugoAlmanaxDate).
			Msgf("DofusDude API returns 404 NOT_FOUND for specific date, continuing with nil almanax...")
	}

	return dodugoAlmanax, nil
}

func (service *Impl) GetAlmanaxByEffect(ctx context.Context, effect, language string,
) (*dodugo.AlmanaxEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoAlmanax *dodugo.AlmanaxEntry
	var dodugoAlmanaxOccurrences []dodugo.AlmanaxEntry
	key := buildKey(almanax, effect, language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoAlmanaxOccurrences) {
		resp, r, err := service.dofusDudeClient.AlmanaxAPI.
			GetAlmanaxRange(ctx, language).
			FilterBonusType(effect).
			RangeSize(constants.DofusDudeAlmanaxSizeLimit).
			Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
	}

	if len(dodugoAlmanaxOccurrences) > 0 {
		dodugoAlmanax = &dodugoAlmanaxOccurrences[0]
	}

	return dodugoAlmanax, nil
}

func (service *Impl) GetAlmanaxByRange(ctx context.Context, daysDuration int32, language string,
) ([]dodugo.AlmanaxEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.httpTimeout)
	defer cancel()

	var dodugoAlmanax []dodugo.AlmanaxEntry
	dodugoAlmanaxDate := time.Now().Format(constants.DofusDudeAlmanaxDateFormat)
	key := buildKey(almanaxRange, fmt.Sprintf("%v_%v", dodugoAlmanaxDate, daysDuration),
		language, constants.GetEncyclopediasSource().Name)
	if !service.getElementFromCache(ctx, key, &dodugoAlmanax) {
		resp, r, err := service.dofusDudeClient.AlmanaxAPI.
			GetAlmanaxRange(ctx, language).
			RangeSize(daysDuration).
			Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			return nil, err
		}
		defer r.Body.Close()
		service.putElementToCache(ctx, key, resp)
		dodugoAlmanax = resp
	}

	return dodugoAlmanax, nil
}
