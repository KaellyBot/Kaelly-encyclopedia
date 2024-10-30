package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
)

type Repository interface {
	GetAlmanaxes() ([]entities.Almanax, error)
	Save(almanax entities.Almanax) error
}

type Impl struct {
	db databases.MySQLConnection
}
