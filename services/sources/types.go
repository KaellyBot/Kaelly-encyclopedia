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

type GameEventHandler func()

type Service interface {
	GetItemType(itemType string) amqp.ItemType
	GetIngredientType(ingredientType string) amqp.IngredientType

	SearchAnyItems(ctx context.Context, query, lg string) ([]dodugo.GetGameSearch200ResponseInner, error)
	SearchEquipments(ctx context.Context, query, lg string) ([]dodugo.ItemListEntry, error)
	SearchMounts(ctx context.Context, query, lg string) ([]dodugo.MountListEntry, error)
	SearchSets(ctx context.Context, query, lg string) ([]dodugo.SetListEntry, error)
	SearchAlmanaxEffects(ctx context.Context, query, lg string) ([]dodugo.GetMetaAlmanaxBonuses200ResponseInner, error)

	GetConsumableByID(ctx context.Context, consumableID int32, lg string) (*dodugo.Resource, error)
	GetEquipmentByID(ctx context.Context, equipmentID int32, lg string) (*dodugo.Weapon, error)
	GetMountByID(ctx context.Context, mountID int32, lg string) (*dodugo.Mount, error)
	GetQuestItemByID(ctx context.Context, questItemID int32, lg string) (*dodugo.Resource, error)
	GetResourceByID(ctx context.Context, resourceID int32, lg string) (*dodugo.Resource, error)
	GetSetByID(ctx context.Context, setID int32, lg string) (*dodugo.EquipmentSet, error)
	GetSets(ctx context.Context) ([]dodugo.SetListEntry, error)

	GetEquipmentByQuery(ctx context.Context, query, lg string) (*dodugo.Weapon, error)
	GetMountByQuery(ctx context.Context, query, lg string) (*dodugo.Mount, error)
	GetSetByQuery(ctx context.Context, query, lg string) (*dodugo.EquipmentSet, error)

	GetAlmanaxByDate(ctx context.Context, date time.Time, language string) (*dodugo.AlmanaxEntry, error)
	GetAlmanaxByEffect(ctx context.Context, effect, language string) (*dodugo.AlmanaxEntry, error)
	GetAlmanaxByRange(ctx context.Context, daysDuration int32, language string) ([]dodugo.AlmanaxEntry, error)

	ListenGameEvent(handler GameEventHandler)
}

type Impl struct {
	eventHandlers   []GameEventHandler
	dofusDudeClient *dodugo.APIClient
	storeService    stores.Service
	httpTimeout     time.Duration
	itemTypes       map[string]amqp.ItemType
	ingredientTypes map[string]amqp.IngredientType
}
