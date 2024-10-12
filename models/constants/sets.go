package constants

import "image"

const (
	MinimumSetBonusItems = 2
)

func SetShieldPoint() image.Point {
	return image.Pt(5, 5)
}

func SetHatPoint() image.Point {
	return image.Pt(210, 5)
}

func SetWeaponPoint() image.Point {
	return image.Pt(415, 5)
}
func SetRing1Point() image.Point {
	return image.Pt(5, 210)
}

func SetAmuletPoint() image.Point {
	return image.Pt(210, 210)
}

func SetRing2Point() image.Point {
	return image.Pt(415, 210)
}

func SetCapePoint() image.Point {
	return image.Pt(5, 415)
}

func SetBeltPoint() image.Point {
	return image.Pt(210, 415)
}

func SetPetPoint() image.Point {
	return image.Pt(5, 625)
}

func SetBootsPoint() image.Point {
	return image.Pt(210, 625)
}
