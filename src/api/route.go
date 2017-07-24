package api

import (
	"octlink/mirage/src/utils"

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

	LoadApiConfig(BaseDir)

	//router.LoadHTMLGlob(BaseDir + "frontend/apitest/templates/*.html")
	router.LoadHTMLFiles(BaseDir+"frontend/apitest/templates/apitest.html",
		BaseDir+"frontend/ng/index.html")

	//router.LoadHTMLGlob(BaseDir + "frontend/ng/index.html")

	router.Static("/static", BaseDir+"frontend/apitest/static")
	router.Static("/ng", BaseDir+"frontend/ng")
	router.Static("/app", BaseDir+"frontend/ng/app")

	router.GET("/apitest", api.LoadApiTestPage)

	router.GET("/", api.LoadNgPage)

	router.GET("/api/", api.ApiTest)
	router.POST("/api/", api.ApiDispatch)
	router.POST("/ng/api/", api.ApiDispatch)

	return router
}
