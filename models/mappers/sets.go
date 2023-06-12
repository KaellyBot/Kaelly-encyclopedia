package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
)

func MapSetList(dodugoSets []dodugo.SetListEntry) []*amqp.EncyclopediaSetListAnswer_Set {
	sets := make([]*amqp.EncyclopediaSetListAnswer_Set, 0)

	for _, set := range dodugoSets {
		sets = append(sets, &amqp.EncyclopediaSetListAnswer_Set{
			Id:   fmt.Sprintf("%v", set.AnkamaId),
			Name: *set.Name,
		})
	}

	return sets
}
