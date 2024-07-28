package constants

import amqp "github.com/kaellybot/kaelly-amqp"

type Ingredient struct {
	ID   string
	Name string
	Type amqp.IngredientType
}
