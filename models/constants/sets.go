package constants

import (
	"image"

	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	MinimumSetBonusItems = 2
)

func GetSetPoints() map[amqp.EquipmentType][]image.Point {
	return map[amqp.EquipmentType][]image.Point{
		amqp.EquipmentType_HAT:    {image.Pt(210, 5)},
		amqp.EquipmentType_CLOAK:  {image.Pt(5, 415)},
		amqp.EquipmentType_AMULET: {image.Pt(210, 210)},
		amqp.EquipmentType_RING:   {image.Pt(5, 210), image.Pt(415, 210)},
		amqp.EquipmentType_BELT:   {image.Pt(210, 415)},
		amqp.EquipmentType_BOOT:   {image.Pt(210, 625)},
		amqp.EquipmentType_WEAPON: {image.Pt(415, 5)},
		amqp.EquipmentType_SHIELD: {image.Pt(5, 5)},
		amqp.EquipmentType_PET:    {image.Pt(5, 625)},
	}
}
