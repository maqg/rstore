package image

import "octlink/rstore/utils/octlog"

var logger *octlog.LogConfig

func InitLog(level int) {
	logger = octlog.InitLogConfig("image.log", level)
}

type Image struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	MediaType  string `json:"mediaType"`
	CreateTime int64  `json:"createTime"`
	LastSync   int64  `json:"lastSync"`
	Desc       string `json:"desc"`
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
