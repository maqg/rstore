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
	Arch        string `json:"arch"`
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
	octlog.Warn("image (%s:%s) deleted\n", image.Name, image.Id)
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

	images := GetAllImages("", "", "")
	for _, image := range images {
		if image.Id == id {
			return &image
		}
	}

	octlog.Error("image of %S not exist", id)

	return nil
}

func GetAllImages(account string, mediaType string, keyword string) []Image {

	imagePath := configuration.GetConfig().RootDirectory + "/" + IMAGESTORE_FILE
	octlog.Debug("find image path[%s]\n", imagePath)

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("open file " + imagePath + "error")
		return make([]Image, 0)
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)

	imageList := make([]Image, 0)
	err = json.Unmarshal(data, &imageList)
	if err != nil {
		octlog.Warn("Transfer json bytes error %s\n", err)
		return make([]Image, 0)
	}

	images := make([]Image, 0)
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

func ContainAccount(accounts []string, id string) bool {
	for _, account := range accounts {
		if account == id {
			return true
		}
	}
	return false
}

func GetAccountList() []string {

	images := GetAllImages("", "", "")
	accounts := make([]string, 0)
	for _, image := range images {
		if image.Account != "" && !ContainAccount(accounts, image.Account) {
			accounts = append(accounts, image.Account)
		}
	}

	return accounts
}
