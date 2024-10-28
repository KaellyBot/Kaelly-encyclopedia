package almanaxes

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/almanaxes"
)

type Service interface {
	GetAlmanaxesByEffect(dofusDudeEffectID string) []entities.Almanax
}

type Impl struct {
	almanaxes  map[string][]entities.Almanax
	repository repository.Repository
}
