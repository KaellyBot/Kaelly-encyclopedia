package sets

import (
	"github.com/kaellybot/kaelly-encyclopedia/models/entities"
	repository "github.com/kaellybot/kaelly-encyclopedia/repositories/sets"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/kaellybot/kaelly-encyclopedia/services/sources"
)

const (
	imgurUploadURL = "https://api.imgur.com/3/image"
)

type imgurResponse struct {
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

type Service interface {
	GetSetByDofusDude(ID int32) (entities.Set, bool)
}

type Impl struct {
	sets             map[int32]entities.Set
	sourceService    sources.Service
	equipmentService equipments.Service
	repository       repository.Repository
}
