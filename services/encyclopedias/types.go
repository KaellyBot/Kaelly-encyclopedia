package encyclopedias

import (
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	requestQueueName   = "encyclopedias-requests"
	requestsRoutingkey = "requests.encyclopedias"
	answersRoutingkey  = "answers.encyclopedias"
)

type Service interface {
	Consume() error
}

type Impl struct {
	// TODO is this right element to use?
	dofusdudeClient dodugo.APIClient
	broker          amqp.MessageBroker
	httpTimeout     time.Duration
}
