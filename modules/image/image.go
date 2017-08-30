package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
	"os"
	"strings"
)

var logger *octlog.LogConfig

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("image.log", level)
}

const (
	// ImageStoreFile for image basic info store file
	ImageStoreFile = "imagestore_info.json"
)

const (
	// ImageStatusReady for ready state
	ImageStatusReady = "ready"

	// ImageStatusDownloading for downloading state
	ImageStatusDownloading = "downloading"

	//ImageStatusError for error status
	ImageStatusError = "error"
)

const (
	// ImageStateEnabled for image state of Enabled
	ImageStateEnabled = "Enabled"

	// ImageStateDisabled for image state disabled
	ImageStateDisabled = "Disabled"
)

// Image for Image sturcture
type Image struct {
	ID          string `json:"uuid"`
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
	URL         string `json:"url"`
	GuestOsType string `json:"guestOsType"` // Guest OS Type
	Arch        string `json:"arch"`
	Platform    string `json:"platform"`
	Format      string `json:"format"`
	System      bool   `json:"isSystem"`
	Account     string `json:"account"`
	InstallPath string `json:"installPath"` // rstore://iamgeid/blobsum
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// GetImageCount to return image count by condition
func GetImageCount(account string, mediaType string, keyword string) int {
	all := GetAllImages(account, mediaType, keyword)
	return len(all)
}

// Brief to return brief info for image
func (image *Image) Brief() map[string]string {
	return map[string]string{
		"id":   image.ID,
		"name": image.Name,
	}
}

// Update to update image
func (image *Image) Update() int {
	all := GetAllImages("", "", "")

	for i, im := range all {
		if im.ID == image.ID {
			all[i] = *image
		}
	}

	WriteImageConfig(all)

	return 0
}

// Add for image, after image added,
// installpath, diskSize, virtualSize, Status, md5sum need update after manifest installed
func (image *Image) Add() int {
	all := GetAllImages("", "", "")
	all = append(all, *image)
	WriteImageConfig(all)
	return 0
}

// Delete for image
func (image *Image) Delete() int {
	octlog.Warn("image (%s:%s) deleted\n", image.Name, image.ID)

	all := GetAllImages("", "", "")
	len := len(all)

	for i, im := range all {
		if im.ID == image.ID {
			if i == 0 {
				all = all[1:len]
			} else {
				all = append(all[0:i], all[i+1:len]...)
			}
		}
	}

	WriteImageConfig(all)

	return 0
}

// FindImageByName find image by name
func FindImageByName(name string) *Image {

	image := new(Image)

	image.Name = "testimage"
	image.ID = "fffffffffffffff"

	octlog.Debug("id %s, name :%s", image.ID, image.Name)

	return image
}

// FindImage by ID
func FindImage(id string) *Image {

	images := GetAllImages("", "", "")
	for _, image := range images {
		if image.ID == id {
			return &image
		}
	}

	octlog.Error("image of %S not exist", id)

	return nil
}

// WriteImageConfig to write all images to image store file
func WriteImageConfig(images []Image) error {

	imagePath := configuration.GetConfig().RootDirectory + "/" + ImageStoreFile
	utils.Remove(imagePath)

	fd, err := os.Create(imagePath)
	if err != nil {
		octlog.Error("create file of %s error\n", imagePath)
		return err
	}

	_, err = fd.Write(utils.JSON2Bytes(images))
	if err != nil {
		// roll back
		return err
	}

	return nil
}

// GetAllImages by condition
func GetAllImages(account string, mediaType string, keyword string) []Image {

	imagePath := configuration.GetConfig().RootDirectory + "/" + ImageStoreFile
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

// ContainAccount check account cantaination
func ContainAccount(accounts []string, id string) bool {
	for _, account := range accounts {
		if account == id {
			return true
		}
	}
	return false
}

// GetAccountList return account list for image
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
