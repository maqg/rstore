package api

import (
	"fmt"
	"octlink/rstore/modules/image"
	"octlink/rstore/utils/merrors"
	"octlink/rstore/utils/octlog"
	"octlink/rstore/utils/uuid"
)

func APIAddImage(paras *ApiParas) *ApiResponse {

	resp := new(ApiResponse)

	newImage := image.FindImageByName(paras.InParas.Paras["image"].(string))
	if newImage != nil {
		logger.Errorf("image %s already exist\n", newImage.Name)
		resp.Error = merrors.ERR_SEGMENT_ALREADY_EXIST
		return resp
	}

	newImage = new(image.Image)
	newImage.Id = uuid.Generate().Simple()
	newImage.Name = paras.InParas.Paras["image"].(string)
	newImage.Desc = paras.InParas.Paras["desc"].(string)

	resp.Error = newImage.Add()

	return resp
}

func APIShowImage(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	imageId := paras.InParas.Paras["id"].(string)
	temp := image.FindImage(imageId)

	if temp == nil {
		resp.Error = merrors.ERR_SEGMENT_NOT_EXIST
		resp.ErrorLog = fmt.Sprintf("user %s not found", imageId)
		return resp
	}

	resp.Data = temp

	octlog.Debug("found User %s", temp.Name)

	return resp
}

func APIUpdateImage(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	id := paras.InParas.Paras["id"].(string)

	ac := image.FindImage(id)
	if ac == nil {
		resp.Error = merrors.ERR_USER_NOT_EXIST
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

func APIDeleteImageByAccount(paras *ApiParas) *ApiResponse {

	octlog.Debug("running in APIDeleteImage\n")

	resp := new(ApiResponse)

	image := image.FindImage(paras.InParas.Paras["id"].(string))
	if image == nil {
		resp.Error = merrors.ERR_SEGMENT_NOT_EXIST
		return resp
	}

	resp.Error = image.Delete()

	return resp
}

func APIDeleteImage(paras *ApiParas) *ApiResponse {

	octlog.Debug("running in APIDeleteImage\n")

	resp := new(ApiResponse)

	image := image.FindImage(paras.InParas.Paras["id"].(string))
	if image == nil {
		resp.Error = merrors.ERR_SEGMENT_NOT_EXIST
		return resp
	}

	resp.Error = image.Delete()

	return resp
}

func APIShowAllImages(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

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

func APIShowAccountList(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	octlog.Debug("running in APIShowAllImage\n")
	accounts := make([]string, 0)
	resp.Data = accounts

	return resp
}
