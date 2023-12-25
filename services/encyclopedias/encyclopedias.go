package encyclopedias

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBroker, sourceService sources.Service,
	equipmentService equipments.Service) *Impl {
	service := Impl{
		sourceService:    sourceService,
		equipmentService: equipmentService,
		broker:           broker,
	}
	service.getListByFunc = map[amqp.EncyclopediaListRequest_Type]getListFunc{
		amqp.EncyclopediaListRequest_ITEM:           service.getItemList,
		amqp.EncyclopediaListRequest_SET:            service.getSetList,
		amqp.EncyclopediaListRequest_ALMANAX_EFFECT: service.getAlmanaxEffectList,
	}

	service.getItemByFuncs = map[amqp.ItemType]getItemFuncs{
		amqp.ItemType_ANY_ITEM: {
			GetItemByID:    service.getItemByID,
			GetItemByQuery: service.getItemByQuery,
		},
		amqp.ItemType_CONSUMABLE: {
			GetItemByID:       service.getConsumableByID,
			GetItemByQuery:    service.getConsumableByQuery,
			GetIngredientByID: service.getConsumableIngredientByID,
		},
		amqp.ItemType_COSMETIC: {
			GetItemByID:    service.getCosmeticByID,
			GetItemByQuery: service.getCosmeticByQuery,
		},
		amqp.ItemType_EQUIPMENT: {
			GetItemByID:       service.getEquipmentByID,
			GetItemByQuery:    service.getEquipmentByQuery,
			GetIngredientByID: service.getEquipmentIngredientByID,
		},
		amqp.ItemType_MOUNT: {
			GetItemByID:    service.getMountByID,
			GetItemByQuery: service.getMountByQuery,
		},
		amqp.ItemType_QUEST_ITEM: {
			GetItemByID:       service.getQuestItemByID,
			GetItemByQuery:    service.getQuestItemByQuery,
			GetIngredientByID: service.getQuestItemIngredientByID,
		},
		amqp.ItemType_RESOURCE: {
			GetItemByID:       service.getResourceByID,
			GetItemByQuery:    service.getResourceByQuery,
			GetIngredientByID: service.getResourceIngredientByID,
		},
		amqp.ItemType_SET: {
			GetItemByID:    service.getSetByID,
			GetItemByQuery: service.getSetByQuery,
		},
	}
	return &service
}

func GetBinding() amqp.Binding {
	return amqp.Binding{
		Exchange:   amqp.ExchangeRequest,
		RoutingKey: requestsRoutingkey,
		Queue:      requestQueueName,
	}
}

func (service *Impl) Consume() error {
	log.Info().Msgf("Consuming encyclopedia requests...")
	return service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(ctx context.Context,
	message *amqp.RabbitMQMessage, correlationID string) {
	//exhaustive:ignore Don't need to be exhaustive here since they will be handled by default case
	switch message.Type {
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_REQUEST:
		service.almanaxRequest(ctx, message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_RESOURCE_REQUEST:
		service.almanaxResourceRequest(ctx, message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_EFFECT_REQUEST:
		service.almanaxEffectRequest(ctx, message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_REQUEST:
		service.listRequest(ctx, message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_REQUEST:
		service.itemRequest(ctx, message, correlationID)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Type not recognized, request ignored")
	}
}
