package api

import (
	"fmt"
	"octlink/rstore/modules/config"
	"octlink/rstore/modules/image"
	"octlink/rstore/utils"
	"octlink/rstore/utils/merrors"
	"octlink/rstore/utils/octlog"
	"octlink/rstore/utils/uuid"
)

// AddImage to add image by API
func AddImage(paras *Paras) *Response {

	id := paras.Get("id")
	if id != "" {
		im := image.GetImage(id)
		if im != nil {
			return &Response{
				Error:    merrors.ErrSegmentAlreadyExist,
				ErrorLog: "User " + id + "Already Exist",
			}
		}
	} else {
		id = uuid.Generate().Simple()
	}

	im := &image.Image{
		ID:          id,
		Arch:        paras.Get("arch"),
		Platform:    paras.Get("platform"),
		GuestOsType: paras.Get("guestOsType"),
		Name:        paras.Get("name"),
		Desc:        paras.Get("desc"),
		MediaType:   paras.Get("mediaType"),
		Format:      paras.Get("format"),
		AccountID:   paras.Get("accountId"),
		CreateTime:  utils.CurrentTimeStr(),
		System:      paras.GetBoolean("isSystem"),
		URL:         paras.Get("url1"),
		Status:      config.ImageStatusDownloading,
		State:       image.ImageStateEnabled,
		LastSync:    utils.CurrentTimeStr(),
	}

	data, err := im.Add()
	return &Response{
		Data:  data,
		Error: err,
	}
}

// ShowImage by api
func ShowImage(paras *Paras) *Response {

	resp := new(Response)
	imageID := paras.InParas.Paras["id"].(string)

	temp := image.GetImage(imageID)
	if temp == nil {
		resp.Error = merrors.ErrSegmentNotExist
		resp.ErrorLog = fmt.Sprintf("user %s not found", imageID)
		return resp
	}

	resp.Data = temp

	return resp
}

// UpdateImage to update image by api
func UpdateImage(paras *Paras) *Response {
	resp := new(Response)

	id := paras.Get("id")
	im := image.GetImage(id)
	if im == nil {
		resp.Error = merrors.ErrUserNotExist
		resp.ErrorLog = "User " + id + "Not Exist"
		return resp
	}

	// Update Image here
	im.Platform = paras.Get("platform")
	im.GuestOsType = paras.Get("guestOsType")
	im.Name = paras.Get("name")
	im.AccountID = paras.Get("accountId")
	if paras.Get("arch") != "" {
		im.Arch = paras.Get("arch")
	}
	if paras.Get("mediaType") != "" {
		im.MediaType = paras.Get("mediaType")
	}
	if paras.Get("format") != "" {
		im.Format = paras.Get("format")
	}

	ret := im.Update()
	if ret != 0 {
		resp.Error = ret
		octlog.Error("update image of %s error\n", im.ID)
		return resp
	}

	return resp
}

// DeleteImageByAccount to delete image by account
func DeleteImageByAccount(paras *Paras) *Response {

	octlog.Debug("running in APIDeleteImage\n")

	resp := new(Response)

	images := image.GetAllImages(paras.InParas.Paras["accountId"].(string), "", "")
	for _, image := range images {
		err := image.Delete()
		if err != 0 {
			resp.Error = err
			octlog.Error("Errored when deleting image of %s\n", image.Name)
			return resp
		}
	}

	return resp
}

// DeleteImage to delete image
func DeleteImage(paras *Paras) *Response {

	octlog.Debug("running in APIDeleteImage\n")

	resp := new(Response)

	image := image.GetImage(paras.Get("id"))
	if image == nil {
		resp.Error = merrors.ErrSegmentNotExist
		return resp
	}

	resp.Error = image.Delete()

	return resp
}

// ShowAllImages to display all images by condition
func ShowAllImages(paras *Paras) *Response {
	resp := new(Response)

	octlog.Debug("running in APIShowAllImage\n")

	imageList := image.GetAllImages(paras.Get("accountId"),
		paras.Get("mediaType"), paras.Get("keyword"))
	resp.Data = imageList
	return resp
}

// ShowAccountList of this rstore server
func ShowAccountList(paras *Paras) *Response {
	return &Response{
		Data: image.GetAccountList(),
	}
}
