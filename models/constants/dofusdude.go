package constants

import amqp "github.com/kaellybot/kaelly-amqp"

const (
	DofusDudeGame            = "dofus2"
	DofusDudeDefaultLanguage = "en"
	DofusDudeLimit           = 25
)

func GetLanguages() map[amqp.Language]string {
	return map[amqp.Language]string{
		amqp.Language_ANY: DofusDudeDefaultLanguage,
		amqp.Language_FR:  "fr",
		amqp.Language_EN:  "en",
		amqp.Language_ES:  "es",
		amqp.Language_DE:  "de",
	}
}
