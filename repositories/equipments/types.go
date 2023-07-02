package equipments

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
)

type Repository interface {
	GetEquipmentTypes() ([]entities.EquipmentType, error)
}

type Impl struct {
	db databases.MySQLConnection
}
