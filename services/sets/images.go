package sets

import (
	"bytes"
	"fmt"
	"image"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
)

func (service *Impl) buildSetImage(items []*dodugo.Weapon) (*bytes.Buffer, error) {
	slotGrid, errSlotGrid := imaging.Open("resources/slot-grid.png")
	if errSlotGrid != nil {
		return nil, errSlotGrid
	}

	slot, errSlot := imaging.Open("resources/slot.png")
	if errSlot != nil {
		return nil, errSlot
	}

	defaultItem, errDefault := imaging.Open("resources/default-item.png")
	if errDefault != nil {
		return nil, errDefault
	}

	var ringNumber int
	for _, item := range items {
		itemImage := getImageFromItem(item, defaultItem)
		equipType, typeFound := service.equipmentService.GetTypeByDofusDude(*item.GetType().Id)
		if !typeFound {
			return nil, fmt.Errorf("item %v type not recognized: %v",
				item.GetAnkamaId(), *item.GetType().Id)
		}

		index := 0
		if equipType.ID == amqp.EquipmentType_RING {
			index = index + ringNumber
			ringNumber = ringNumber + 1
		}

		points, pointFound := constants.GetSetPoints()[equipType.ID]
		if !pointFound {
			return nil, fmt.Errorf("item %v type have not equivalent point: %v",
				item.GetAnkamaId(), *item.GetType().Id)
		}

		slotGrid = appendImage(slotGrid, slot, itemImage, points[index])
	}

	buf, errBuf := imageToBuffer(slotGrid)
	if errBuf != nil {
		return nil, errBuf
	}

	return buf, nil
}

func appendImage(itemGrid, slot, item image.Image,
	point image.Point) image.Image {
	itemSlot := imaging.Overlay(slot, item, image.Point{0, 0}, 1)
	return imaging.Overlay(itemGrid, itemSlot, point, 1)
}

func getImageFromURL(url string) (image.Image, error) {
	resp, errGet := http.Get(url)
	if errGet != nil {
		return nil, errGet
	}
	defer resp.Body.Close()

	image, errDecode := imaging.Decode(resp.Body)
	if errDecode != nil {
		return nil, errDecode
	}

	return image, nil
}

func getImageFromItem(item *dodugo.Weapon,
	defaultItem image.Image) image.Image {
	if item.GetImageUrls().Sd.IsSet() {
		itemImage, errGetImg := getImageFromURL(*item.GetImageUrls().Sd.Get())
		if errGetImg != nil {
			log.Warn().Err(errGetImg).
				Str(constants.LogAnkamaID, fmt.Sprintf("%v", item.GetAnkamaId())).
				Msgf("Cannot retrieve item SD icon with DofusDude, continuing with default one")
			return defaultItem
		}

		return itemImage
	}

	log.Warn().
		Str(constants.LogAnkamaID, fmt.Sprintf("%v", item.GetAnkamaId())).
		Msgf("Item SD icon not set with DofusDude, continuing with default one")
	return defaultItem
}

func imageToBuffer(img image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := imaging.Encode(buf, img, imaging.PNG)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
