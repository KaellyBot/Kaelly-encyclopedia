package encyclopedias

import (
	"context"
	"errors"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
)

const (
	requestQueueName   = "encyclopedias-requests"
	requestsRoutingkey = "requests.encyclopedias"
	answersRoutingkey  = "answers.encyclopedias"
)

var (
	errBadRequestMessage = errors.New("message request could not be satisfied")
)

type getItemListFunc func(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemListAnswer, error)
type getItemByIDFunc func(ctx context.Context, ID int32, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error)
type getItemByQueryFunc func(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error)
type getIngredientByIDFunc func(ctx context.Context, ID int32, correlationID,
	lg string) (*constants.Ingredient, error)

type getItemFuncs struct {
	GetItemByID       getItemByIDFunc
	GetItemByQuery    getItemByQueryFunc
	GetIngredientByID getIngredientByIDFunc
}

type Service interface {
	Consume() error
}

type Impl struct {
	sourceService     sources.Service
	equipmentService  equipments.Service
	broker            amqp.MessageBroker
	getItemByFuncs    map[amqp.ItemType]getItemFuncs
	getItemListByFunc map[amqp.EncyclopediaItemListRequest_Type]getItemListFunc
}
