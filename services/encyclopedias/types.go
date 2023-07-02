package encyclopedias

import (
	"errors"
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
)

const (
	requestQueueName   = "encyclopedias-requests"
	requestsRoutingkey = "requests.encyclopedias"
	answersRoutingkey  = "answers.encyclopedias"
)

type objectType string

const (
	almanax objectType = "almanax"
	item    objectType = "items"
	set     objectType = "sets"
)

var (
	errNotFound = errors.New("cannot find the desired resource")
)

type Service interface {
	Consume() error
}

type Impl struct {
	dofusDudeClient  *dodugo.APIClient
	equipmentService equipments.Service
	storeService     stores.Service
	broker           amqp.MessageBroker
	httpTimeout      time.Duration
}
