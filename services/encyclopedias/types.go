package encyclopedias

import (
	"context"
	"errors"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sets"
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

type getListFunc func(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaListAnswer, error)
type getItemByIDFunc func(ctx context.Context, ID int32, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error)
type getItemByQueryFunc func(ctx context.Context, query, correlationID,
	lg string) (*amqp.EncyclopediaItemAnswer, error)
type getIngredientByIDFunc func(ctx context.Context, ID int32, correlationID,
	lg string) (*constants.Ingredient, error)

type getItemFuncs struct {
	GetItemByID    getItemByIDFunc
	GetItemByQuery getItemByQueryFunc
}

type Service interface {
	Consume() error
}

type Impl struct {
	sourceService        sources.Service
	equipmentService     equipments.Service
	setService           sets.Service
	broker               amqp.MessageBroker
	getListByFunc        map[amqp.EncyclopediaListRequest_Type]getListFunc
	getItemByFuncs       map[amqp.ItemType]getItemFuncs
	getIngredientByFuncs map[amqp.ItemType]getIngredientByIDFunc
}
