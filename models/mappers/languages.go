package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
)

func MapLanguage(lg amqp.Language) string {
	language, found := constants.GetLanguages()[lg]
	if !found {
		language = constants.DofusDudeDefaultLanguage
	}
	return language
}
