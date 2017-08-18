package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"octlink/rstore/configuration"
	"octlink/rstore/utils/octlog"
	"os"
	"strings"
)

var logger *octlog.LogConfig

func InitLog(level int) {
	logger = octlog.InitLogConfig("image.log", level)
}

const (
	IMAGESTORE_FILE = "imagestore_info.json"
)

type Image struct {
	Id          string `json:"uuid"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Status      string `json:"status"`
	MediaType   string `json:"mediaType"`
	CreateTime  string `json:"createTime"`
	LastSync    string `json:"lastSync"`
	Desc        string `json:"description"`
	DiskSize    int64  `json:"diskSize"`
	VirtualSize int64  `json:"virtualSize"`
	Md5Sum      string `json:"md5sum"`
	Url         string `json:"url"`
	Type        string `json:"type"`
	Platform    string `json:"platform"`
	Format      string `json:"format"`
	System      bool   `json:"system"`
	Account     string `json:"account"`
	InstallPath string `json:"installPath"`
}

func GetImageCount() int {
	return 10
}

func (image *Image) Brief() map[string]string {
	return map[string]string{
		"id":   image.Id,
		"name": image.Name,
	}
}

func (image *Image) Update() int {
	return 0
}

func (image *Image) Add() int {
	return 0
}

func (image *Image) Delete() int {
	octlog.Debug("image deleted\n")
	return 0
}

func FindImageByName(name string) *Image {

	image := new(Image)

	image.Name = "testimage"
	image.Id = "fffffffffffffff"

	octlog.Debug("id %s, name :%s", image.Id, image.Name)

	return image
}

func FindImage(id string) *Image {

	image := new(Image)

	image.Name = "testimage"
	image.Id = "fffffffffffffff"

	octlog.Debug("id %s, name :%s", image.Id, image.Name)

	return image
}

func GetAllImages(account string, mediaType string, keyword string) []Image {

	imagePath := configuration.GetConfig().RootDirectory + "/" + IMAGESTORE_FILE

	octlog.Debug("find image path[%s]\n", imagePath)

	imageList := make([]Image, 0)
	images := make([]Image, 0)

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("open file " + imagePath + "error")
		return nil
	}

	data, err := ioutil.ReadFile(imagePath)
	defer file.Close()

	err = json.Unmarshal(data, &imageList)
	if err != nil {
		octlog.Warn("Transfer json bytes error %s\n", err)
		return nil
	}

	for _, image := range imageList {

		// filter account
		if account != "" && image.Account != account {
			continue
		}

		// filter mediaType
		if mediaType != "" && image.MediaType != mediaType {
			continue
		}

		// filter keyword
		if keyword != "" && !strings.Contains(image.Name, keyword) {
			continue
		}

		images = append(images, image)
	}

	return images
}
