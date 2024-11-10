package encyclopedias

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/almanaxes"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBroker, sourceService sources.Service,
	almanaxService almanaxes.Service, equipmentService equipments.Service,
	setService sets.Service) *Impl {
	service := Impl{
		sourceService:    sourceService,
		almanaxService:   almanaxService,
		equipmentService: equipmentService,
		setService:       setService,
		broker:           broker,
	}

	service.getListByFunc = map[amqp.EncyclopediaListRequest_Type]getListFunc{
		amqp.EncyclopediaListRequest_UNKNOWN:        nil,
		amqp.EncyclopediaListRequest_ITEM:           service.getItemList,
		amqp.EncyclopediaListRequest_SET:            service.getSetList,
		amqp.EncyclopediaListRequest_ALMANAX_EFFECT: service.getAlmanaxEffectList,
	}

	//nolint:exhaustive // Not all types can be requested at this moment.
	service.getItemByFuncs = map[amqp.ItemType]getItemFuncs{
		amqp.ItemType_ANY_ITEM_TYPE: {
			GetItemByID:    service.getItemByID,
			GetItemByQuery: service.getItemByQuery,
		},
		amqp.ItemType_COSMETIC_TYPE: {
			GetItemByID:    service.getCosmeticByID,
			GetItemByQuery: service.getCosmeticByQuery,
		},
		amqp.ItemType_EQUIPMENT_TYPE: {
			GetItemByID:    service.getEquipmentByID,
			GetItemByQuery: service.getEquipmentByQuery,
		},
		amqp.ItemType_MOUNT_TYPE: {
			GetItemByID:    service.getMountByID,
			GetItemByQuery: service.getMountByQuery,
		},
		amqp.ItemType_SET_TYPE: {
			GetItemByID:    service.getSetByID,
			GetItemByQuery: service.getSetByQuery,
		},
	}

	//nolint:exhaustive // Ingredient types possibility is exhaustive here.
	service.getIngredientByFuncs = map[amqp.ItemType]getIngredientByIDFunc{
		amqp.ItemType_ANY_ITEM_TYPE:   nil,
		amqp.ItemType_CONSUMABLE_TYPE: service.getConsumableIngredientByID,
		amqp.ItemType_EQUIPMENT_TYPE:  service.getEquipmentIngredientByID,
		amqp.ItemType_QUEST_ITEM_TYPE: service.getQuestItemIngredientByID,
		amqp.ItemType_RESOURCE_TYPE:   service.getResourceIngredientByID,
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

func (service *Impl) consume(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	//exhaustive:ignore Don't need to be exhaustive here since they will be handled by default case
	switch message.Type {
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_REQUEST:
		service.almanaxRequest(ctx, message)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_RESOURCE_REQUEST:
		service.almanaxResourceRequest(ctx, message)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_EFFECT_REQUEST:
		service.almanaxEffectRequest(ctx, message)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_REQUEST:
		service.listRequest(ctx, message)
	case amqp.RabbitMQMessage_ENCYCLOPEDIA_ITEM_REQUEST:
		service.itemRequest(ctx, message)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Msgf("Type not recognized, request ignored")
	}
}

func (service *Impl) replyWithSuceededAnswer(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	err := service.broker.Reply(message, ctx.CorrelationID, ctx.ReplyTo)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogReplyTo, ctx.ReplyTo).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) replyWithFailedAnswer(ctx amqp.Context, messageType amqp.RabbitMQMessage_Type,
	language amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     messageType,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: language,
	}

	err := service.broker.Reply(&message, ctx.CorrelationID, ctx.ReplyTo)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogReplyTo, ctx.ReplyTo).
			Msgf("Cannot publish via broker, request ignored")
	}
}
