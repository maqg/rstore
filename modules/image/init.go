package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"octlink/rstore/modules/config"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
	"os"
	"os/signal"
	"syscall"
)

const (
	// MaxImagesCount for max images count
	MaxImagesCount = 1000
)

// GImages for all image loaded from config
var GImages []*Image

// GImagesMap Global Images Map
var GImagesMap map[string]*Image

// GImagesIsoMap iso map list
var GImagesIsoMap map[string]*Image

// GImagesRootTemplateMap for root template map
var GImagesRootTemplateMap map[string]*Image

// GImagesDataTemplateMap for data volume template map
var GImagesDataTemplateMap map[string]*Image

func loadImagesFromConfig() error {

	imagePath := configuration.GetConfig().RootDirectory + "/" + ImageStoreFile
	if !utils.IsFileExist(imagePath) {
		octlog.Error("file of %s not exist\n", imagePath)
		return fmt.Errorf("file of %s not exist", imagePath)
	}

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

func zeroImages() {
	GImages = make([]*Image, 0)
	GImagesMap = make(map[string]*Image, MaxImagesCount)
	GImagesIsoMap = make(map[string]*Image, MaxImagesCount)
	GImagesRootTemplateMap = make(map[string]*Image, MaxImagesCount)
	GImagesDataTemplateMap = make(map[string]*Image, MaxImagesCount)
}

// ReloadImages for images reloading
func ReloadImages() error {

	// zero images firstly
	zeroImages()

	err := loadImagesFromConfig()
	if err != nil {
		octlog.Error("load images error [%s]\n", err)
		return nil
	}

	for _, im := range GImages {

		// global map
		GImagesMap[im.ID] = im

		switch im.MediaType {
		case config.ImageTypeRootTemplate:
			GImagesRootTemplateMap[im.ID] = im
			break

		case config.ImageTypeDataVolume:
			GImagesDataTemplateMap[im.ID] = im

		case config.ImageTypeIso:
			GImagesIsoMap[im.ID] = im
			break
		}
	}

	return nil
}

// WriteImages to write all images to image store file
func WriteImages() error {

	imagePath := configuration.GetConfig().RootDirectory + "/" + ImageStoreFile
	if utils.IsFileExist(imagePath) {
		os.Rename(imagePath, imagePath+"."+utils.CurrentTimeSimple())
	}

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

// ReloadSignal for image reload signal handler
func ReloadSignal() {
	c := make(chan os.Signal, 1)
	// 10 for SIGUSR1, 12 for SIGUSR2
	signal.Notify(c, syscall.Signal(10), syscall.Signal(12))
	for {
		s := <-c
		fmt.Print("Got signal:", s)
		fmt.Println(", now reload images from config file")
		ReloadImages()
	}
}
