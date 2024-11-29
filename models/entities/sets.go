package entities

import amqp "github.com/kaellybot/kaelly-amqp"

type Set struct {
	DofusDudeID int32     `gorm:"primaryKey"`
	Game        amqp.Game `gorm:"primaryKey"`
	Icon        string
}
