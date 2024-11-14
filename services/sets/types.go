package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/news"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
)

const (
	setBaseURL = "https://raw.githubusercontent.com/KaellyBot/Kaelly-cdn/refs/heads/main/KaellyBot/sets/%v.webp"
)

type Service interface {
	GetSetByDofusDude(ID int32) (entities.Set, bool)
}

type Impl struct {
	sets             map[int32]entities.Set
	newsService      news.Service
	sourceService    sources.Service
	equipmentService equipments.Service
	repository       repository.Repository
}
