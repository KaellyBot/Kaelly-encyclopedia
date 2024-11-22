package games

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetGameVersion(id amqp.Game) (entities.GameVersion, error) {
	var gameVersion entities.GameVersion
	response := repo.db.GetDB().
		Where(&entities.GameVersion{ID: id}).
		First(&gameVersion)
	return gameVersion, response.Error
}

func (repo *Impl) Save(gameVersion entities.GameVersion) error {
	return repo.db.GetDB().Save(&gameVersion).Error
}
