package mappers

import (
	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
)

func MapMount(item *dodugo.Mount) *amqp.EncyclopediaItemAnswer {

	// TODO

	return &amqp.EncyclopediaItemAnswer{
		Type:  amqp.ItemType_MOUNT,
		Mount: &amqp.EncyclopediaItemAnswer_Mount{},
		Source: &amqp.Source{
			Name: constants.GetEncyclopediasSource().Name,
			Icon: constants.GetEncyclopediasSource().Icon,
			Url:  constants.GetEncyclopediasSource().URL,
		},
	}
}
