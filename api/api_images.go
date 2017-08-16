package api

import (
	"fmt"
	"octlink/rstore/modules/image"
	"octlink/rstore/utils"
	"octlink/rstore/utils/merrors"
	"octlink/rstore/utils/octlog"
	"octlink/rstore/utils/uuid"
)

func APIAddImage(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	newImage := image.FindImageByName(paras.Db, paras.InParas.Paras["image"].(string))
	if newImage != nil {
		logger.Errorf("image %s already exist\n", newImage.Name)
		resp.Error = merrors.ERR_SEGMENT_ALREADY_EXIST
		return resp
	}

	newImage = new(image.Image)
	newImage.Id = uuid.Generate().Simple()
	newImage.Name = paras.InParas.Paras["image"].(string)
	newImage.Type = utils.ParasInt(paras.InParas.Paras["type"])
	newImage.Email = paras.InParas.Paras["email"].(string)
	newImage.PhoneNumber = paras.InParas.Paras["phoneNumber"].(string)
	newImage.Password = paras.InParas.Paras["password"].(string)
	newImage.Desc = paras.InParas.Paras["desc"].(string)

	resp.Error = newImage.Add(paras.Db)

	return resp
}

func APIShowImage(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	imageId := paras.InParas.Paras["id"].(string)
	temp := image.FindImage(paras.Db, imageId)

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

	ac := image.FindImage(paras.Db, id)
	if ac == nil {
		resp.Error = merrors.ERR_USER_NOT_EXIST
		resp.ErrorLog = "User " + id + "Not Exist"
		return resp
	}

	ac.Email = paras.InParas.Paras["email"].(string)
	ac.PhoneNumber = paras.InParas.Paras["phoneNumber"].(string)
	ac.Desc = paras.InParas.Paras["desc"].(string)

	ret := ac.Update(paras.Db)
	if ret != 0 {
		resp.Error = ret
		return resp
	}

	return resp
}

func APIDeleteImage(paras *ApiParas) *ApiResponse {

	octlog.Debug("running in APIDeleteImage\n")

	resp := new(ApiResponse)

	image := image.FindImage(paras.Db,
		paras.InParas.Paras["id"].(string))
	if image == nil {
		resp.Error = merrors.ERR_SEGMENT_NOT_EXIST
		return resp
	}

	resp.Error = image.Delete(paras.Db)

	return resp
}

func APIShowAllImage(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	octlog.Debug("running in APIShowAllImage\n")

	offset := utils.ParasInt(paras.InParas.Paras["start"])
	limit := utils.ParasInt(paras.InParas.Paras["limit"])

	rows, err := paras.Db.Query("SELECT ID,U_Name,U_State,U_Type,U_Email,U_PhoneNumber,"+
		"U_Description,U_CreateTime,U_LastLogin,U_LastSync "+
		"FROM tb_image LIMIT ?,?", offset, limit)
	if err != nil {
		logger.Errorf("query image list error %s\n", err.Error())
		resp.Error = merrors.ERR_DB_ERR
		return resp
	}
	defer rows.Close()

	imageList := make([]image.Image, 0)

	for rows.Next() {
		var image image.Image
		err = rows.Scan(&image.Id, &image.Name, &image.State,
			&image.Type, &image.Email, &image.PhoneNumber, &image.Desc,
			&image.CreateTime, &image.LastLogin, &image.LastSync)
		if err == nil {
			logger.Debugf("query result: %s:%s\n", image.Id,
				image.Name, image.State, image.Type)
		} else {
			logger.Errorf("query image list error %s\n", err.Error())
		}
		imageList = append(imageList, image)
	}

	count := image.GetImageCount(paras.Db)

	result := make(map[string]interface{}, 3)
	result["total"] = count
	result["count"] = len(imageList)
	result["data"] = imageList

	resp.Data = result

	return resp
}
