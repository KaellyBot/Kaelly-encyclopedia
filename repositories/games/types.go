package games

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
)

type Repository interface {
	GetGameVersion(id amqp.Game) (entities.GameVersion, error)
	Save(entity entities.GameVersion) error
}

type Impl struct {
	db databases.MySQLConnection
}
