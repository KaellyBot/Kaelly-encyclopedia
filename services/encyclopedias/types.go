package encyclopedias

import (
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
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

type Service interface {
	Consume() error
}

type Impl struct {
	dofusDudeClient *dodugo.APIClient
	storeService    stores.Service
	broker          amqp.MessageBroker
	httpTimeout     time.Duration
}
