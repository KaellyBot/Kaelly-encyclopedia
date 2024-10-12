package sets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/disintegration/imaging"
	"github.com/dofusdude/dodugo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-encyclopedia/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func (service *Impl) buildSetImage(ctx context.Context, items []*dodugo.Weapon,
) (*bytes.Buffer, error) {
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
		itemImage := getImageFromItem(ctx, item, defaultItem)
		equipType, typeFound := service.equipmentService.GetTypeByDofusDude(*item.GetType().Id)
		if !typeFound {
			return nil, fmt.Errorf("item %v type not recognized: %v",
				item.GetAnkamaId(), *item.GetType().Id)
		}

		index := 0
		if equipType.ID == amqp.EquipmentType_RING {
			index += ringNumber
			ringNumber++
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

func getImageFromURL(ctx context.Context, rawURL string) (image.Image, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if errReq != nil {
		return nil, errReq
	}

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer resp.Body.Close()

	image, errDecode := imaging.Decode(resp.Body)
	if errDecode != nil {
		return nil, errDecode
	}

	return image, nil
}

func getImageFromItem(ctx context.Context, item *dodugo.Weapon,
	defaultItem image.Image) image.Image {
	if item.GetImageUrls().Sd.IsSet() {
		itemImage, errGetImg := getImageFromURL(ctx, *item.GetImageUrls().Sd.Get())
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

func uploadImageToImgur(ctx context.Context, buf *bytes.Buffer) (string, error) {
	// Create a multipart writer
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	// Create the file field for the image data
	part, errForm := writer.CreateFormFile("image", "image.png")
	if errForm != nil {
		return "", errForm
	}

	// Write the image data to the file field
	_, errCopy := io.Copy(part, buf)
	if errCopy != nil {
		return "", errCopy
	}

	// Close the multipart writer to finalize the body
	errClose := writer.Close()
	if errClose != nil {
		return "", errClose
	}

	// Create the request
	req, errReq := http.NewRequestWithContext(ctx, http.MethodPost, imgurUploadURL, reqBody)
	if errReq != nil {
		return "", errReq
	}

	// Set necessary headers
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %v",
		viper.GetString(constants.ImgurClientID)))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return "", errDo
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("imgur API error, status code: %d", resp.StatusCode)
	}

	var imgurResponse imgurResponse
	decoder := json.NewDecoder(resp.Body)
	if errDecode := decoder.Decode(&imgurResponse); errDecode != nil {
		return "", errDecode
	}

	if !imgurResponse.Success {
		return "", fmt.Errorf("imgur API error, inner status code: %d", imgurResponse.Status)
	}

	return imgurResponse.Data.Link, nil
}
