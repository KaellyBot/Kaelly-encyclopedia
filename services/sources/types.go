package sources

import (
	"context"
	"errors"
	"time"

	"github.com/dofusdude/dodugo"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
)

type objectType string

const (
	almanax objectType = "almanax"
	item    objectType = "items"
	set     objectType = "sets"
)

var (
	errNotFound = errors.New("cannot find the desired resource")
)

type Service interface {
	SearchAnyItems(ctx context.Context, query, lg string) ([]dodugo.ItemsListEntryTyped, error)
	SearchEquipments(ctx context.Context, query, lg string) ([]dodugo.ItemListEntry, error)
	SearchSets(ctx context.Context, query, lg string) ([]dodugo.SetListEntry, error)

	GetEquipmentByID(ctx context.Context, itemID int32, lg string) (*dodugo.Weapon, error)
	GetSetByID(ctx context.Context, setID int32, lg string) (*dodugo.EquipmentSet, error)

	GetEquipmentByQuery(ctx context.Context, query, lg string) (*dodugo.Weapon, error)
	GetSetByQuery(ctx context.Context, query, lg string) (*dodugo.EquipmentSet, error)
}

type Impl struct {
	dofusDudeClient *dodugo.APIClient
	storeService    stores.Service
	httpTimeout     time.Duration
}
