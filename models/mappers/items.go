package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
)

func MapItemList(dodugoItems []dodugo.ItemsListEntryTyped) []*amqp.EncyclopediaItemListAnswer_Item {
	items := make([]*amqp.EncyclopediaItemListAnswer_Item, 0)

	for _, item := range dodugoItems {
		items = append(items, &amqp.EncyclopediaItemListAnswer_Item{
			Id:   fmt.Sprintf("%v", item.AnkamaId),
			Name: *item.Name,
		})
	}

	return items
}
