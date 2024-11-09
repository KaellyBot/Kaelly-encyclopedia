package entities

import amqp "github.com/kaellybot/kaelly-amqp"

type EquipmentType struct {
	EquipmentID amqp.EquipmentType `gorm:"primaryKey"`
	ItemID      amqp.ItemType      `gorm:"primaryKey"`
	DofusDudeID int32              `gorm:"primaryKey"`
}
