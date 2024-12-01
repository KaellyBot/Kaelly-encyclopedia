package mappers

import (
	"fmt"

	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/kaellybot/kaelly-encyclopedia/services/equipments"
	"github.com/rs/zerolog/log"
)

func MapEquipment(item *dodugo.Weapon, ingredientItems map[int32]*constants.Ingredient,
	equipmentService equipments.Service) *amqp.EncyclopediaItemAnswer {
	var set *amqp.EncyclopediaItemAnswer_Equipment_SetFamily
	if item.HasParentSet() {
		parentSet := item.GetParentSet()
		set = &amqp.EncyclopediaItemAnswer_Equipment_SetFamily{
			Id:   fmt.Sprintf("%v", parentSet.GetId()),
			Name: parentSet.GetName(),
		}
	}

	weaponEffects := make([]*amqp.EncyclopediaItemAnswer_Effect, 0)
	effects := make([]*amqp.EncyclopediaItemAnswer_Effect, 0)
	for _, effect := range item.GetEffects() {
		amqpEffect := &amqp.EncyclopediaItemAnswer_Effect{
			Id:    fmt.Sprintf("%v", *effect.GetType().Id),
			Label: effect.GetFormatted(),
		}

		if effect.GetType().IsActive != nil && *effect.GetType().IsActive {
			weaponEffects = append(weaponEffects, amqpEffect)
		} else {
			effects = append(effects, amqpEffect)
		}
	}

	var recipe *amqp.EncyclopediaItemAnswer_Recipe
	if len(item.GetRecipe()) > 0 {
		ingredients := make([]*amqp.EncyclopediaItemAnswer_Recipe_Ingredient, 0)
		for _, recipeEntry := range item.GetRecipe() {
			formattedItemIDString := fmt.Sprintf("%v", recipeEntry.GetItemAnkamaId())
			ingredient, found := ingredientItems[recipeEntry.GetItemAnkamaId()]
			if !found {
				log.Warn().
					Str(constants.LogAnkamaID, formattedItemIDString).
					Msgf("Cannot build entire recipe (missing ingredient), continuing with degraded mode")
				ingredient = &constants.Ingredient{
					Name: formattedItemIDString,
					Type: amqp.ItemType_ANY_ITEM_TYPE,
				}
			}

			ingredients = append(ingredients, &amqp.EncyclopediaItemAnswer_Recipe_Ingredient{
				Id:       fmt.Sprintf("%v", recipeEntry.GetItemAnkamaId()),
				Name:     ingredient.Name,
				Quantity: int64(recipeEntry.GetQuantity()),
				Type:     ingredient.Type,
			})
		}

		recipe = &amqp.EncyclopediaItemAnswer_Recipe{
			Ingredients: ingredients,
		}
	}

	icon := item.GetImageUrls().Icon
	if item.GetImageUrls().Hq.IsSet() {
		icon = item.GetImageUrls().Hq.Get()
	}

	equipmentType := mapEquipmentType(item.GetType(), equipmentService)

	var characteristics *amqp.EncyclopediaItemAnswer_Equipment_Characteristics
	if item.GetIsWeapon() {
		characteristics = &amqp.EncyclopediaItemAnswer_Equipment_Characteristics{
			Cost:           int64(item.GetApCost()),
			MinRange:       int64(item.Range.GetMin()),
			MaxRange:       int64(item.Range.GetMax()),
			MaxCastPerTurn: int64(item.GetMaxCastPerTurn()),
			CriticalRate:   int64(item.GetCriticalHitProbability()),
			CriticalBonus:  int64(item.GetCriticalHitBonus()),
			// TODO area
		}
	}

	return &amqp.EncyclopediaItemAnswer{
		Type: amqp.ItemType_EQUIPMENT_TYPE,
		Equipment: &amqp.EncyclopediaItemAnswer_Equipment{
			Id:          fmt.Sprintf("%v", item.GetAnkamaId()),
			Name:        item.GetName(),
			Description: item.GetDescription(),
			Type: &amqp.EncyclopediaItemAnswer_Equipment_Type{
				ItemType:       equipmentType.ItemID,
				EquipmentType:  equipmentType.EquipmentID,
				EquipmentLabel: *item.GetType().Name,
			},
			Icon:            *icon,
			Level:           int64(item.GetLevel()),
			Pods:            int64(item.GetPods()),
			Set:             set,
			Characteristics: characteristics,
			WeaponEffects:   weaponEffects,
			Effects:         effects,
			Conditions:      mapNullableConditions(item.Conditions),
			Recipe:          recipe,
		},
		Source: constants.GetDofusDudeSource(),
	}
}

func mapNullableConditions(conditions dodugo.NullableConditionNode,
) *amqp.EncyclopediaItemAnswer_Conditions {
	if !conditions.IsSet() {
		return nil
	}

	return mapConditions(conditions.Get())
}

func mapConditions(conditions *dodugo.ConditionNode,
) *amqp.EncyclopediaItemAnswer_Conditions {
	if conditions == nil {
		return nil
	}

	leaf := conditions.ConditionLeaf
	if leaf != nil {
		return &amqp.EncyclopediaItemAnswer_Conditions{
			Relation: amqp.EncyclopediaItemAnswer_Conditions_NONE,
			Condition: &amqp.EncyclopediaItemAnswer_Conditions_Condition{
				Operator: leaf.Condition.GetOperator(),
				Value:    int64(leaf.Condition.GetIntValue()),
				Element: &amqp.EncyclopediaItemAnswer_Conditions_Condition_Element{
					Id:   fmt.Sprintf("%v", leaf.Condition.Element.GetId()),
					Name: leaf.Condition.Element.GetName(),
				},
			},
		}
	}

	innerConditions := conditions.ConditionRelation
	if innerConditions != nil {
		var relation amqp.EncyclopediaItemAnswer_Conditions_Relation
		switch innerConditions.GetRelation() {
		case "or":
			relation = amqp.EncyclopediaItemAnswer_Conditions_OR
		case "and":
			relation = amqp.EncyclopediaItemAnswer_Conditions_AND
		default:
			log.Warn().
				Msgf("Cannot determine properly item condition relation '%v', using '%v' by default",
					innerConditions.GetRelation(), amqp.EncyclopediaItemAnswer_Conditions_NONE)
			relation = amqp.EncyclopediaItemAnswer_Conditions_NONE
		}

		children := make([]*amqp.EncyclopediaItemAnswer_Conditions, 0)
		for _, child := range innerConditions.GetChildren() {
			node := child
			children = append(children, mapConditions(node))
		}

		return &amqp.EncyclopediaItemAnswer_Conditions{
			Relation: relation,
			Children: children,
		}
	}

	return nil
}
