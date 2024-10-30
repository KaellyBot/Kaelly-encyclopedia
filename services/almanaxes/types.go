package almanaxes

import (
	"time"

	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/almanaxes"
)

type Service interface {
	GetDatesByAlmanaxEffect(dofusDudeEffectID string) []time.Time
}

type Impl struct {
	almanaxes  map[string][]entities.Almanax
	repository repository.Repository
}
