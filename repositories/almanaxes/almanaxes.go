package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetAlmanaxes() ([]entities.Almanax, error) {
	var almanaxes []entities.Almanax
	response := repo.db.GetDB().
		Model(&entities.Almanax{}).
		Find(&almanaxes)
	return almanaxes, response.Error
}
