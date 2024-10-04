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
		amqp.EncyclopediaListRequest_UNKNOWN:        nil,
		amqp.EncyclopediaListRequest_ITEM:           service.getItemList,
		amqp.EncyclopediaListRequest_SET:            service.getSetList,
		amqp.EncyclopediaListRequest_ALMANAX_EFFECT: service.getAlmanaxEffectList,
	}

	service.getItemByFuncs = map[amqp.ItemType]getItemFuncs{
		amqp.ItemType_ANY_ITEM: {
			GetItemByID:    service.getItemByID,
			GetItemByQuery: service.getItemByQuery,
		},
		amqp.ItemType_EQUIPMENT: {
			GetItemByID:    service.getEquipmentByID,
			GetItemByQuery: service.getEquipmentByQuery,
		},
		amqp.ItemType_SET: {
			GetItemByID:    service.getSetByID,
			GetItemByQuery: service.getSetByQuery,
		},
	}

	service.getIngredientByFuncs = map[amqp.IngredientType]getIngredientByIDFunc{
		amqp.IngredientType_ANY_INGREDIENT:       nil,
		amqp.IngredientType_CONSUMABLE:           service.getConsumableIngredientByID,
		amqp.IngredientType_EQUIPMENT_INGREDIENT: service.getEquipmentIngredientByID,
		amqp.IngredientType_QUEST_ITEM:           service.getQuestItemIngredientByID,
		amqp.IngredientType_RESOURCE:             service.getResourceIngredientByID,
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
