package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	"github.com/kaellybot/kaelly-encyclopedia/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetSets() ([]entities.Set, error) {
	var sets []entities.Set
	response := repo.db.GetDB().Model(&entities.Set{}).Find(&sets)
	return sets, response.Error
}

func (repo *Impl) Save(set entities.Set) error {
	return repo.db.GetDB().Save(&set).Error
}
