package image

import (
	"fmt"
	"octlink/rstore/modules/config"
	"octlink/rstore/modules/task"
	"octlink/rstore/utils"
	"octlink/rstore/utils/merrors"
	"octlink/rstore/utils/octlog"
	"octlink/rstore/utils/uuid"
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
	WriteImages()
	return 0
}

// UpdateImageCallback when image download OK, update its info
func UpdateImageCallback(imageID string, diskSize int64, virtualSize int64,
	blobsum string, status string) error {

	im := GetImage(imageID)
	if im != nil {
		octlog.Warn("Got image of %s\n", imageID)
		im.Status = config.ImageStatusReady
		im.DiskSize = diskSize
		im.VirtualSize = virtualSize
		im.Md5Sum = blobsum
		im.InstallPath = fmt.Sprintf("rstore://%s/%s", im.ID, im.Md5Sum)
		im.Status = status

		WriteImages()

		return nil
	}

	return fmt.Errorf("Image of %s not exist", imageID)
}

// Add for image, after image added,
// installpath, diskSize, virtualSize, Status, md5sum need update after manifest installed
func (image *Image) Add() (string, int) {

	GImages = append(GImages, image)
	GImagesMap[image.ID] = image

	WriteImages()

	if image.URL != "" {
		t := new(task.Task)
		t.ID = uuid.Generate().Simple()
		t.URL = image.URL
		t.CreateTime = utils.CurrentTimeStr()
		t.ImageName = image.ID
		t.Status = task.TaskStatusNew
		t.Callback = UpdateImageCallback
		t.AddAndRun(UpdateImageCallback)
		return t.ID, merrors.ErrSuccess
	}

	return "", merrors.ErrSuccess
}

// Delete for image
func (image *Image) Delete() int {

	octlog.Warn("image (%s:%s) deleted\n", image.Name, image.ID)

	len := len(GImages)

	for i, im := range GImages {
		if im.ID == image.ID {
			if i == 0 {
				GImages = GImages[1:len]
			} else {
				GImages = append(GImages[0:i], GImages[i+1:len]...)
			}
			delete(GImagesMap, im.ID)
		}
	}

	WriteImages()

	return 0
}

// GetImage by ID
func GetImage(id string) *Image {
	return GImagesMap[id]
}

// GetAllImages by condition
func GetAllImages(account string, mediaType string, keyword string) []*Image {

	images := make([]*Image, 0)

	for _, image := range GImages {

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
