package api

import (
	"fmt"
	"octlink/rstore/modules/image"
	"octlink/rstore/utils/merrors"
	"octlink/rstore/utils/octlog"
	"octlink/rstore/utils/uuid"
)

// AddImage to add image by API
func AddImage(paras *Paras) *Response {

	resp := new(Response)

	newImage := new(image.Image)
	newImage.ID = uuid.Generate().Simple()
	newImage.Name = paras.InParas.Paras["image"].(string)
	newImage.Desc = paras.InParas.Paras["desc"].(string)

	resp.Error = newImage.Add()

	return resp
}

// ShowImage by api
func ShowImage(paras *Paras) *Response {

	resp := new(Response)
	imageID := paras.InParas.Paras["id"].(string)

	temp := image.FindImage(imageID)
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

	id := paras.InParas.Paras["id"].(string)

	ac := image.FindImage(id)
	if ac == nil {
		resp.Error = merrors.ErrUserNotExist
		resp.ErrorLog = "User " + id + "Not Exist"
		return resp
	}

	ac.Desc = paras.InParas.Paras["desc"].(string)

	ret := ac.Update()
	if ret != 0 {
		resp.Error = ret
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

	image := image.FindImage(paras.InParas.Paras["id"].(string))
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

	imageList := image.GetAllImages(paras.InParas.Paras["accountId"].(string),
		paras.InParas.Paras["mediaType"].(string), paras.InParas.Paras["keyword"].(string))

	result := make(map[string]interface{}, 3)
	result["total"] = 0
	result["count"] = len(imageList)
	result["data"] = imageList

	resp.Data = result

	return resp
}

// ShowAccountList of this rstore server
func ShowAccountList(paras *Paras) *Response {

	octlog.Debug("running in APIShowAllImage\n")

	resp := new(Response)

	resp.Data = image.GetAccountList()

	return resp
}
