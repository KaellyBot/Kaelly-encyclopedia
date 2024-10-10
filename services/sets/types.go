package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
)

type Service interface {
	GetSetByDofusDude(ID int32) (entities.Set, bool)
}

type Impl struct {
	sets       map[int32]entities.Set
	repository repository.Repository
}
