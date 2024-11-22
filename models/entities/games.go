package entities

import amqp "github.com/kaellybot/kaelly-amqp"

type GameVersion struct {
	ID      amqp.Game `gorm:"primaryKey"`
	Version string
}
