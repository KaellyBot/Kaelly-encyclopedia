package constants

import (
	"image"

	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	MinimumSetBonusItems = 2

	setItemMarginPx = 5
	setItemSizePx   = 200
	setFirstCell    = setItemMarginPx
	setSecondCell   = setFirstCell + setItemSizePx + setItemMarginPx
	setThirdCell    = setSecondCell + setItemSizePx + setItemMarginPx
	setFourthCell   = setThirdCell + setItemSizePx + setItemMarginPx
)

//nolint:exhaustive // No other types needed.
func GetSetPoints() map[amqp.EquipmentType][]image.Point {
	return map[amqp.EquipmentType][]image.Point{
		amqp.EquipmentType_SHIELD: {image.Pt(setFirstCell, setFirstCell)},
		amqp.EquipmentType_HAT:    {image.Pt(setSecondCell, setFirstCell)},
		amqp.EquipmentType_WEAPON: {image.Pt(setThirdCell, setFirstCell)},
		amqp.EquipmentType_AMULET: {image.Pt(setSecondCell, setSecondCell)},
		amqp.EquipmentType_RING: {
			image.Pt(setFirstCell, setSecondCell),
			image.Pt(setThirdCell, setSecondCell),
		},
		amqp.EquipmentType_CLOAK: {image.Pt(setFirstCell, setThirdCell)},
		amqp.EquipmentType_BELT:  {image.Pt(setSecondCell, setThirdCell)},
		amqp.EquipmentType_PET:   {image.Pt(setFirstCell, setFourthCell)},
		amqp.EquipmentType_BOOT:  {image.Pt(setSecondCell, setFourthCell)},
	}
}
