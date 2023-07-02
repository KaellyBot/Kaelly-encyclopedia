package entities

import amqp "github.com/kaellybot/kaelly-amqp"

type EquipmentType struct {
	ID          amqp.EquipmentType `gorm:"primaryKey"`
	DofusDudeID int32              `gorm:"primaryKey"`
	DebugName   string
}
