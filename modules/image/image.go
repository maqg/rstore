package image

import (
	"fmt"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/config"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/modules/task"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
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
	ID          string `json:"id"`
	Name        string `json:"imageName"`
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
	AccountID   string `json:"accountId"`
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

// UpdateImageCallback when image download OK, update its info,
// if Image not exist, just create a new one and update it.
func UpdateImageCallback(imageID string, diskSize int64, virtualSize int64,
	blobsum string, status string) error {

	im := GetImage(imageID)
	if im == nil {
		im = new(Image)

		// default type to RootVolumeTemplate
		im.State = ImageStateEnabled
		im.Name = imageID
		im.ID = imageID
		im.MediaType = config.ImageTypeRootTemplate
		im.CreateTime = utils.CurrentTimeStr()
		im.LastSync = utils.CurrentTimeStr()
		im.GuestOsType = "unknown"
		im.Arch = "amd64"
		im.Platform = "Linux"
		im.Format = "qcow2"

		appendImage(im)
	}
	im.Status = config.ImageStatusReady
	im.DiskSize = diskSize
	im.VirtualSize = virtualSize
	im.Md5Sum = blobsum
	im.InstallPath = fmt.Sprintf("rstore://%s/%s", im.ID, im.Md5Sum)
	im.Status = status

	WriteImages()

	return nil
}

// Add for image, after image added,
// installpath, diskSize, virtualSize, Status, md5sum need update after manifest installed
func (image *Image) Add() (string, int) {

	appendImage(image)

	WriteImages()

	if image.URL != "" {
		t := new(task.Task)
		t.ID = uuid.Generate().Simple()
		t.URL = image.URL
		t.CreateTime = utils.CurrentTimeStr()
		t.ImageName = image.ID
		t.Status = task.TaskStatusNew
		t.AddAndRun(UpdateImageCallback)
		return t.ID, merrors.ErrSuccess
	}

	return "", merrors.ErrSuccess
}

func (image *Image) removeManifest() {
	m := manifest.GetManifest(image.ID, utils.GetDigestStr(image.ID))
	if m != nil {
		m.Delete()
		logger.Warnf("deleted manifest (%s:%s)\n", image.Name, m.BlobSum)
	}
}

func isBlobsManifestUsed(imageID string, blobsum string) bool {

	for _, im := range GImages {
		if im.Md5Sum == blobsum && im.ID != imageID {
			return true
		}
	}

	return false
}

func (image *Image) removeBlobsManifest() {

	bm := blobsmanifest.GetBlobsManifest(image.Md5Sum)
	if bm == nil {
		logger.Errorf("blobs-manifest of %s not exist\n", image.Md5Sum)
		return
	}

	if !isBlobsManifestUsed(image.ID, image.Md5Sum) {
		bm.Delete()
		logger.Warnf("deleted blobs-manifest (%s:%s)\n", image.ID, image.Md5Sum)
	}
}

func (image *Image) removeManifestDir() {
	baseDir := configuration.GetConfig().RootDirectory + fmt.Sprintf(manifest.ImageManifestDirProto, image.ID)
	utils.RemoveDir(baseDir)
}

// append image to list
func appendImage(im *Image) {

	// append to Global image list
	GImages = append(GImages, im)

	// to Images map
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

// remove image from list
func removeImage(im *Image) {

	// remove from images map
	delete(GImagesMap, im.ID)

	// remove from root images map
	delete(GImagesRootTemplateMap, im.ID)

	// remove from data images map
	delete(GImagesDataTemplateMap, im.ID)

	// remove from iso images map
	delete(GImagesIsoMap, im.ID)
}

// Delete for image
func (image *Image) Delete() int {

	logger.Warnf("now to delete image (%s:%s)\n", image.Name, image.ID)

	len := len(GImages)

	for i, im := range GImages {
		if im.ID == image.ID {
			if i == 0 {
				GImages = GImages[1:len]
			} else {
				GImages = append(GImages[0:i], GImages[i+1:len]...)
			}

			removeImage(im)

			// To remove manifest firstly
			im.removeManifest()

			// To remove base manifest directory
			im.removeManifestDir()

			// then remove all blobs
			im.removeBlobsManifest()

			logger.Warnf("deleted image (%s:%s)\n", image.Name, image.ID)
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
		if account != "" && image.AccountID != account {
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
		if image.AccountID != "" && !ContainAccount(accounts, image.AccountID) {
			accounts = append(accounts, image.AccountID)
		}
	}

	return accounts
}
