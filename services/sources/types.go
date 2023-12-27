package sources

import (
	"context"
	"errors"
	"time"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/services/stores"
)

type objectType string

const (
	almanax       objectType = "almanax"
	almanaxRange  objectType = "almanaxRange"
	almanaxEffect objectType = "almanaxEffect"
	item          objectType = "items"
	set           objectType = "sets"
)

var (
	ErrNotFound = errors.New("cannot find the desired resource")
)

type Service interface {
	GetItemType(itemType string) amqp.ItemType

	SearchAnyItems(ctx context.Context, query, lg string) ([]dodugo.ItemsListEntryTyped, error)
	SearchConsumables(ctx context.Context, query, lg string) ([]dodugo.ItemListEntry, error)
	SearchCosmetics(ctx context.Context, query, lg string) ([]dodugo.ItemListEntry, error)
	SearchEquipments(ctx context.Context, query, lg string) ([]dodugo.ItemListEntry, error)
	SearchMounts(ctx context.Context, query, lg string) ([]dodugo.MountListEntry, error)
	SearchQuestItems(ctx context.Context, query, lg string) ([]dodugo.ItemListEntry, error)
	SearchResources(ctx context.Context, query, lg string) ([]dodugo.ItemListEntry, error)
	SearchSets(ctx context.Context, query, lg string) ([]dodugo.SetListEntry, error)
	SearchAlmanaxEffects(ctx context.Context, query, lg string) ([]dodugo.GetMetaAlmanaxBonuses200ResponseInner, error)

	GetConsumableByID(ctx context.Context, consumableID int32, lg string) (*dodugo.Resource, error)
	GetCosmeticByID(ctx context.Context, cosmeticID int32, lg string) (*dodugo.Cosmetic, error)
	GetEquipmentByID(ctx context.Context, equipmentID int32, lg string) (*dodugo.Weapon, error)
	GetMountByID(ctx context.Context, mountID int32, lg string) (*dodugo.Mount, error)
	GetQuestItemByID(ctx context.Context, questItemID int32, lg string) (*dodugo.Resource, error)
	GetResourceByID(ctx context.Context, resourceID int32, lg string) (*dodugo.Resource, error)
	GetSetByID(ctx context.Context, setID int32, lg string) (*dodugo.EquipmentSet, error)

	GetConsumableByQuery(ctx context.Context, query, lg string) (*dodugo.Resource, error)
	GetCosmeticByQuery(ctx context.Context, query, lg string) (*dodugo.Cosmetic, error)
	GetEquipmentByQuery(ctx context.Context, query, lg string) (*dodugo.Weapon, error)
	GetMountByQuery(ctx context.Context, query, lg string) (*dodugo.Mount, error)
	GetQuestItemByQuery(ctx context.Context, query, lg string) (*dodugo.Resource, error)
	GetResourceByQuery(ctx context.Context, query, lg string) (*dodugo.Resource, error)
	GetSetByQuery(ctx context.Context, query, lg string) (*dodugo.EquipmentSet, error)

	GetAlmanaxByDate(ctx context.Context, date time.Time, language string) (*dodugo.AlmanaxEntry, error)
	GetAlmanaxByRange(ctx context.Context, daysDuration int32, language string) ([]dodugo.AlmanaxEntry, error)
}

type Impl struct {
	dofusDudeClient *dodugo.APIClient
	storeService    stores.Service
	httpTimeout     time.Duration
	itemTypes       map[string]amqp.ItemType
}
