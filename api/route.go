package api

import (
	"octlink/rstore/utils"

	"github.com/gin-gonic/gin"
)

func (api *Api) ApiRouter() *gin.Engine {

	var BaseDir string = ""

	router := gin.New()

	gin.SetMode(gin.ReleaseMode)

	exist := utils.IsFileExist("frontend")
	if !exist {
		BaseDir = "../"
	}

	//LoadApiConfig(BaseDir)

	router.LoadHTMLGlob(BaseDir + "frontend/apitest/templates/*.html")
	router.Static("/static", BaseDir+"frontend/static")

	router.GET("/api/test/", api.LoadApiTestPage)

	router.GET("/api/", api.ApiTest)
	router.POST("/api/", api.ApiDispatch)

	return router
}
