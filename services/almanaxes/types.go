package almanaxes

import (
	"errors"
	"time"

	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/almanaxes"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
)

var (
	errNotFound = errors.New("almanax is not found")
)

type Service interface {
	GetDatesByAlmanaxEffect(dofusDudeEffectID string) []time.Time
}

type Impl struct {
	almanaxes     map[string][]entities.Almanax
	sourceService sources.Service
	repository    repository.Repository
}
