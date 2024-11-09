package constants

import (
	"image"

	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	setItemMarginPx = 5
	setItemSizePx   = 200
	setFirstCell    = setItemMarginPx
	setSecondCell   = setFirstCell + setItemSizePx + setItemMarginPx
	setThirdCell    = setSecondCell + setItemSizePx + setItemMarginPx
	setFourthCell   = setThirdCell + setItemSizePx + setItemMarginPx
)

//nolint:exhaustive // No other types needed.
func GetSetPoints() map[amqp.EquipmentType][]image.Point {
	weaponPoints := []image.Point{image.Pt(setThirdCell, setFirstCell)}
	return map[amqp.EquipmentType][]image.Point{

		amqp.EquipmentType_HAT:     {image.Pt(setSecondCell, setFirstCell)},
		amqp.EquipmentType_CLOAK:   {image.Pt(setFirstCell, setThirdCell)},
		amqp.EquipmentType_AXE:     weaponPoints,
		amqp.EquipmentType_BOW:     weaponPoints,
		amqp.EquipmentType_DAGGER:  weaponPoints,
		amqp.EquipmentType_HAMMER:  weaponPoints,
		amqp.EquipmentType_LANCE:   weaponPoints,
		amqp.EquipmentType_PICKAXE: weaponPoints,
		amqp.EquipmentType_SCYTHE:  weaponPoints,
		amqp.EquipmentType_SHOVEL:  weaponPoints,
		amqp.EquipmentType_STAFF:   weaponPoints,
		amqp.EquipmentType_SWORD:   weaponPoints,
		amqp.EquipmentType_WAND:    weaponPoints,
		amqp.EquipmentType_SHIELD:  {image.Pt(setFirstCell, setFirstCell)},
		amqp.EquipmentType_PET:     {image.Pt(setFirstCell, setFourthCell)},

		amqp.EquipmentType_AMULET: {image.Pt(setSecondCell, setSecondCell)},
		amqp.EquipmentType_RING: {
			image.Pt(setFirstCell, setSecondCell),
			image.Pt(setThirdCell, setSecondCell),
		},
		amqp.EquipmentType_BELT: {image.Pt(setSecondCell, setThirdCell)},
		amqp.EquipmentType_BOOT: {image.Pt(setSecondCell, setFourthCell)},
	}
}
