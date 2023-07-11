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
	service.getItemListByFunc = map[amqp.EncyclopediaItemListRequest_Type]getItemListFunc{
		amqp.EncyclopediaItemListRequest_ANY: service.getItemList,
		amqp.EncyclopediaItemListRequest_SET: service.getSetList,
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
		// TODO do others too
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
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_LIST_REQUEST:
		service.itemListRequest(ctx, message, correlationID)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_REQUEST:
		service.itemRequest(ctx, message, correlationID)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Type not recognized, request ignored")
	}
}
