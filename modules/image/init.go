package image

import (
	"encoding/json"
	"io/ioutil"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
	"os"
)

const (
	// MaxImagesCount for max images count
	MaxImagesCount = 1000
)

// GImages for all image loaded from config
var GImages = make([]Image, MaxImagesCount)

// GImagesMap Global Images Map
var GImagesMap = make(map[string]Image, MaxImagesCount)

func loadImages() error {

	imagePath := configuration.GetConfig().RootDirectory + "/" + ImageStoreFile
	file, err := os.Open(imagePath)
	if err != nil {
		octlog.Error("open image store file " + imagePath + "error\n")
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	err = json.Unmarshal(data, &GImages)
	if err != nil {
		octlog.Warn("Transfer json bytes error %s\n", err)
		return err
	}

	return nil
}

// ReloadImages for images reloading
func ReloadImages() error {

	err := loadImages()
	if err != nil {
		octlog.Error("load images error [%s]\n", err)
		return nil
	}

	for _, im := range GImages {
		GImagesMap[im.ID] = im
	}

	return nil
}

// WriteImages to write all images to image store file
func WriteImages() error {

	imagePath := configuration.GetConfig().RootDirectory + "/" + ImageStoreFile
	utils.Remove(imagePath)

	fd, err := os.Create(imagePath)
	if err != nil {
		octlog.Error("create file of %s error\n", imagePath)
		return err
	}

	_, err = fd.Write(utils.JSON2Bytes(GImages))
	if err != nil {
		octlog.Error("Write images to image store file %s error\n", imagePath)
		// roll back
		return err
	}

	return nil
}
